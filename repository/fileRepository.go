package repository

import (
	"feklistova/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// CreateFile creates a new file in the database and returns the created file
func (r *Repository) CreateFile(ctx context.Context, file *models.File, isDownloaded bool) error {
	var fileID int
	var err error
	if isDownloaded {
		err = r.Db.QueryRowContext(ctx, `
			INSERT INTO downloaded_files (shipment_id, filepath, timestamp)
			VALUES ($1, $2, $3)
			RETURNING file_id`,
			file.ShipmentID, file.FilePath, file.Timestamp).Scan(&fileID)
	} else {
		err = r.Db.QueryRowContext(ctx, `
			INSERT INTO model_files (shipment_id, filepath, timestamp)
			VALUES ($1, $2, $3)
			RETURNING file_id`,
			file.ShipmentID, file.FilePath, file.Timestamp).Scan(&fileID)
	}
	if err != nil {
		return err
	}

	file.FileID = fileID
	return nil
}

func (r *Repository) CreateModelFile(ctx context.Context, file *models.File, metrics []models.ModelMetrics) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = r.insertFile(ctx, tx, file)
	if err != nil {
		return err
	}

	fileID := file.FileID

	err = r.saveModelMetrics(ctx, tx, fileID, metrics)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) insertFile(ctx context.Context, tx *sql.Tx, file *models.File) error {
	var fileID int
	err := tx.QueryRowContext(ctx, `
		INSERT INTO model_files (shipment_id, filepath, timestamp)
		VALUES ($1, $2, $3)
		RETURNING file_id`,
		file.ShipmentID, file.FilePath, file.Timestamp).Scan(&fileID)
	if err != nil {
		return fmt.Errorf("failed to insert file: %v", err)
	}

	file.FileID = fileID
	return nil
}

func (r *Repository) UpdateFilePathByID(ctx context.Context, fileID int, newFilePath string, isDownloaded bool) error {
	var tableName string
	if isDownloaded {
		tableName = "downloaded_files"
	} else {
		tableName = "model_files"
	}

	_, err := r.Db.ExecContext(ctx, `
        UPDATE `+tableName+`
        SET filepath = $1
        WHERE file_id = $2`,
		newFilePath, fileID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile deletes a file from the database by its file ID
func (r *Repository) DeleteFile(ctx context.Context, fileID int, isDownloaded bool) error {
	var err error
	if isDownloaded {
		_, err = r.Db.ExecContext(ctx, "DELETE FROM downloaded_files WHERE file_id = $1", fileID)
	} else {
		_, err = r.Db.ExecContext(ctx, "DELETE FROM model_files WHERE file_id = $1", fileID)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetDownloadedFilesByShipmentID(ctx context.Context, shipmentID int) ([]models.File, error) {
	var files []models.File

	downloadedQuery := `
        SELECT file_id, shipment_id, filepath, timestamp FROM downloaded_files WHERE shipment_id = $1
    `
	downloadedRows, err := r.Db.QueryContext(ctx, downloadedQuery, shipmentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query downloaded files")
	}
	defer downloadedRows.Close()

	for downloadedRows.Next() {
		var file models.File
		if err := downloadedRows.Scan(&file.FileID, &file.ShipmentID, &file.FilePath, &file.Timestamp); err != nil {
			return nil, errors.Wrap(err, "failed to scan downloaded file row")
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *Repository) GetUploadedFilesByShipmentID(ctx context.Context, shipmentID int) ([]models.File, error) {
	var files []models.File

	uploadedQuery := `
        SELECT file_id, shipment_id, filepath, timestamp FROM model_files WHERE shipment_id = $1
    `
	uploadedRows, err := r.Db.QueryContext(ctx, uploadedQuery, shipmentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query uploaded files")
	}
	defer uploadedRows.Close()

	for uploadedRows.Next() {
		var file models.File
		if err := uploadedRows.Scan(&file.FileID, &file.ShipmentID, &file.FilePath, &file.Timestamp); err != nil {
			return nil, errors.Wrap(err, "failed to scan uploaded file row")
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *Repository) saveModelMetrics(ctx context.Context, tx *sql.Tx, fileID int, metrics []models.ModelMetrics) error {
	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO model_metrics (file_id, metric_name, metric_value)
        VALUES ($1, $2, $3)
    `)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement")
	}
	defer stmt.Close()

	for _, metric := range metrics {
		_, err := stmt.ExecContext(ctx, fileID, metric.MetricName, metric.MetricValue)
		if err != nil {
			return errors.Wrap(err, "failed to insert metric")
		}
	}

	return nil
}

func (r *Repository) GetMetricsByFileID(ctx context.Context, fileID int) (map[string]float64, error) {
	metrics := make(map[string]float64)

	rows, err := r.Db.QueryContext(ctx, `
        SELECT metric_name, metric_value
        FROM model_metrics
        WHERE file_id = $1
    `, fileID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	for rows.Next() {
		var metricName string
		var metricValue float64
		if err := rows.Scan(&metricName, &metricValue); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		metrics[metricName] = metricValue
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error occurred during iteration")
	}

	return metrics, nil
}
