package main

import (
	"feklistova/models"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func ShipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelType := vars["model_type"]
	if modelType != "reg" && modelType != "class" {
		log.Printf("Restricted http request for shipment: incorrect request, model type %s is not allowed", modelType)
		http.Error(w, "Error retrieving model type", http.StatusBadRequest)
		return
	}

	if !IsAuthorized(r) {
		log.Println("Unauthorized user attempted shipment")
		http.Redirect(w, r, "/users/enter", http.StatusSeeOther)
		return
	}
	userID := GetUserID(r)
	if userID == -1 {
		log.Println("Could not retrieve userID")
		http.Error(w, "Unable to identify user", http.StatusBadRequest)
		return
	}
	log.Printf("User %d requests new shipment", userID)

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("Error parsing shipment form: %v", err)
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Extract individual form fields
	projectName := r.FormValue("project_name")
	algorithm := r.FormValue("model_type")
	targetColumn := r.FormValue("target_column")

	// Handle file upload
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error saving uploaded file: %v", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	filenameParts := strings.Split(fileHeader.Filename, ".")
	fileExtension := ""
	if len(filenameParts) > 1 {
		fileExtension = filenameParts[len(filenameParts)-1]
	}

	log.Printf("Received request from user %d to create project: %s, model: %s, algorithm: %s, target col: %s",
		userID, projectName, modelType, algorithm, targetColumn)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	// Creating shipment
	shipment := &models.Shipment{
		UserID:       userID,
		ProjectName:  projectName,
		ModelType:    modelType,
		Algorithm:    algorithm,
		TargetColumn: targetColumn,
		Status:       "accepted",
		Timestamp:    time.Now(),
	}
	if err := repo.CreateShipment(ctx, shipment); err != nil {
		log.Printf("Error creating shipment: %v", err)
		http.Error(w, "Error creating shipment", http.StatusInternalServerError)
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			shipment.Status = "failed"
		} else if err != nil {
			shipment.Status = "denied"
		} else {
			shipment.Status = "finished"
		}

		ctxFinal, cancelFinal := context.WithTimeout(r.Context(), time.Minute*5)
		defer cancelFinal()

		if err := repo.UpdateShipmentStatus(ctxFinal, shipment); err != nil {
			log.Fatal(err)
		}

		log.Printf("Shipment %d status changed to %s", shipment.ShipmentID, shipment.Status)

		if shipment.Status == "failed" || shipment.Status == "denied" {
			// clearing files
			files, err := repo.GetDownloadedFilesByShipmentID(ctxFinal, shipment.ShipmentID)
			if err != nil {
				log.Printf("Failed to get download files for shipment ID %d: %v", shipment.ShipmentID, err)
				return
			}
			for _, file := range files {
				if err := fileRepo.DeleteDownloadedFile(strconv.Itoa(file.FileID)); err != nil {
					log.Printf("Failed to delete file with ID %d: %v", file.FileID, err)
				}
				if err := repo.DeleteFile(ctxFinal, file.FileID, true); err != nil {
					log.Fatalf("Failed to forget file with ID %d: %v", file.FileID, err)
				}
			}

			files, err = repo.GetUploadedFilesByShipmentID(ctxFinal, shipment.ShipmentID)
			if err != nil {
				log.Printf("Failed to get upload files for shipment ID %d: %v", shipment.ShipmentID, err)
				return
			}
			for _, file := range files {
				if err := fileRepo.DeleteUploadedFile(strconv.Itoa(file.FileID)); err != nil {
					log.Printf("Failed to delete file with ID %d: %v", file.FileID, err)
				}
				if err := repo.DeleteFile(ctxFinal, file.FileID, false); err != nil {
					log.Fatalf("Failed to forget file with ID %d: %v", file.FileID, err)
				}
			}
		}
	}()

	downloadedFile := &models.File{
		ShipmentID: shipment.ShipmentID,
		Timestamp:  time.Now(),
	}
	if err := repo.CreateFile(ctx, downloadedFile, true); err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Error creating downloading file", http.StatusInternalServerError)
		return
	}

	newFileName := fmt.Sprintf("%d.%s", downloadedFile.FileID, fileExtension)
	log.Printf("Saving downloaded file as %s", newFileName)
	if err := fileRepo.SaveDownloadedFile(newFileName, file); err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Error saving downloaded file", http.StatusInternalServerError)
		return
	}

	downloadedFilePath := fileRepo.GetDownloadedFilePath(newFileName)
	if err := repo.UpdateFilePathByID(ctx, downloadedFile.FileID, downloadedFilePath, true); err != nil {
		log.Printf("Error updating file: %v", err)
		http.Error(w, "Error saving downloaded file", http.StatusInternalServerError)
		return
	}

	shipment.Status = "in progress"
	if err := repo.UpdateShipmentStatus(ctx, shipment); err != nil {
		log.Println(err)
		return
	}

	uploadedFilePath := fileRepo.GetUploadedFilePath(strconv.Itoa(downloadedFile.FileID))
	// Start the Python model process
	log.Printf("Running model for shipment %d", shipment.ShipmentID)
	metricsDict, err := pyModel.RunModel(
		modelType,
		algorithm,
		targetColumn,
		downloadedFilePath,
		uploadedFilePath,
	)
	if err != nil {
		log.Printf("Error running python model: %v", err)
		os.Remove(uploadedFilePath)
		http.Error(w, "Error starting Python model", http.StatusInternalServerError)
		return
	}

	metrics := models.ParseMetricsToModelMetrics(downloadedFile.FileID, metricsDict)
	modelOutputFile := &models.File{
		FilePath:   uploadedFilePath,
		ShipmentID: shipment.ShipmentID,
		Timestamp:  time.Now(),
	}
	ctxSaving, cancelSaving := context.WithTimeout(r.Context(), time.Second*5)
	defer cancelSaving()

	if err := repo.CreateModelFile(ctxSaving, modelOutputFile, metrics); err != nil {
		log.Printf("Error creating model file: %v", err)
		http.Error(w, "Error saving model file", http.StatusInternalServerError)
		return
	}
	log.Printf("Model file successfully updated in the database: %s", uploadedFilePath)

	// Redirect to another page where results will be presented
	http.Redirect(w, r, "/shipment/result/"+strconv.Itoa(shipment.ShipmentID), http.StatusFound)
}

func ResultShipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipmentID, err := strconv.Atoi(vars["shipment_id"])
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}
	log.Printf("Results for shipment %d requested", shipmentID)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	shipment, err := repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		log.Printf("Failed to load shipment by ID %d: %v", shipmentID, err)
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	uploadedFiles, err := repo.GetUploadedFilesByShipmentID(ctx, shipmentID)
	if err != nil {
		log.Printf("Failed to get upload files for shipment ID %d: %v", shipmentID, err)
		http.Error(w, "Failed to get upload files", http.StatusInternalServerError)
		return
	}
	if len(uploadedFiles) != 1 {
		log.Printf("Expected length of files to be 1, got %d", len(uploadedFiles))
	}
	metricsDict, err := repo.GetMetricsByFileID(ctx, uploadedFiles[0].FileID)
	if err != nil {
		log.Printf("Failed to retrieve metrics by file ID %d: %v", uploadedFiles[0].FileID, err)
		http.Error(w, "Unable to fetch metrics", http.StatusBadRequest)
		return
	}
	log.Printf("Results for shipment %d requested: %s", shipmentID, uploadedFiles[0].FilePath)
	for name, value := range metricsDict {
		log.Printf("	%s: %.2f\n", name, value)
	}
	if shipment.ModelType == "reg" {
		shipmentRegHandler(w, r, shipmentID, metricsDict)
	} else if shipment.ModelType == "class" {
		shipmentClassHandler(w, r, shipmentID, metricsDict)
	} else {
		log.Printf("Unknown ModelType for shipment ID %d: %s", shipmentID, shipment.ModelType)
		http.Error(w, "Failed to parse ModelType", http.StatusInternalServerError)
		return
	}
}

type RegHandlerMetrics struct {
	ShipmentID string
	R2         string
	MAE        string
}

type ClassHandlerMetrics struct {
	ShipmentID string
	Precision  string
	Recall     string
	F1Score    string
}

func shipmentClassHandler(w http.ResponseWriter, r *http.Request, shipmentID int, metricsDict map[string]float64) {
	r.ParseForm()

	modelParams := ClassHandlerMetrics{
		ShipmentID: fmt.Sprintf("%d", shipmentID),
		Precision:  fmt.Sprintf("%.2f", metricsDict["Precision"]),
		Recall:     fmt.Sprintf("%.2f", metricsDict["Recall"]),
		F1Score:    fmt.Sprintf("%.2f", metricsDict["F1-score"]),
	}

	tmpl, err := template.ParseFiles("web/save_model_class.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, modelParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func shipmentRegHandler(w http.ResponseWriter, r *http.Request, shipmentID int, metricsDict map[string]float64) {
	r.ParseForm()

	modelParams := RegHandlerMetrics{
		ShipmentID: fmt.Sprintf("%d", shipmentID),
		R2:         fmt.Sprintf("%.2f", metricsDict["R2 Score"]),
		MAE:        fmt.Sprintf("%.2f", metricsDict["MAE"]),
	}

	tmpl, err := template.ParseFiles("web/save_model_reg.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, modelParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ShipmentDownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipmentIDStr := vars["shipment_id"]
	shipmentID, err := strconv.Atoi(shipmentIDStr)
	if err != nil {
		log.Printf("Requested to download results with invalid shipment ID %s: %v", shipmentIDStr, err)
		http.Error(w, "Invalid shipment_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	uploadedFiles, err := repo.GetUploadedFilesByShipmentID(ctx, shipmentID)
	if err != nil {
		log.Printf("Failed to get upload files for shipment ID %d: %v", shipmentID, err)
		http.Error(w, "Failed to get upload files", http.StatusInternalServerError)
		return
	}
	if len(uploadedFiles) != 1 {
		log.Printf("Expected length of files to be 1, got %d", len(uploadedFiles))
	}
	filePath := uploadedFiles[0].FilePath
	log.Printf("Results are being sent to download for shipment ID %d: %s", shipmentID, filePath)

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Printf("Results not found for shipment ID %d: %s", shipmentID, filePath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open results for shipment ID %d: %s", shipmentID, filePath)
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=results")
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Failed to send results for shipment ID %d: %s", shipmentID, filePath)
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
