package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server received %s request for %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Serve - является функцией работы сервера
func Serve(ctx context.Context) {
	server := http.Server{Addr: ":8080"}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	http.Handle("/", router)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	assetsHandler := http.FileServer(http.Dir("./web/"))
	prefixHandler := http.StripPrefix("/assets/", assetsHandler)
	http.Handle("/assets/", prefixHandler)
	router.HandleFunc("/home", HomeHandlerTmpl).Methods("GET") // index.html

	// registration
	router.HandleFunc("/users/enter", LoginHandlerTmpl).Methods("GET")       // enter.html
	router.HandleFunc("/users/register", RegisterHandlerTmpl).Methods("GET") // registration.html
	router.HandleFunc("/profile", ProfileHandlerTmpl).Methods("GET")         // profile.html

	router.HandleFunc("/api/users/enter", LoginHandler).Methods("POST") // enter.html
	router.HandleFunc("/api/users/register", RegisterHandler)           // registration.html
	router.HandleFunc("/api/profile", ProfileHandler)                   // profile.html

	// shipment
	router.HandleFunc("/shipment/model_class", ProgressClassHandlerTmpl).Methods("GET") // model_form_class.html
	router.HandleFunc("/shipment/model_reg", ProgressRegHandlerTmpl).Methods("GET")     // model_form_reg.html

	// model_type is expected to be either class or reg
	router.HandleFunc("/api/shipment/progress/{model_type}", ShipmentHandler) // progress_class.html
	router.HandleFunc("/api/shipment/download_results/{shipment_id}", ShipmentDownloadHandler)
	router.HandleFunc("/shipment/result/{shipment_id}", ResultShipmentHandler) // save_model_class.html / save_model_reg.html

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down server")
			err := server.Shutdown(ctx)
			if err != nil {
				panic(err)
			}
			return
		}

	}
}
