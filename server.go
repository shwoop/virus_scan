package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type httpHandler func(w http.ResponseWriter, r *http.Request)

func scanHandlerWrapper(scanner scanFunc) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("scan requested")
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "400: Invalid request, must be multipart form and under 10Mb", http.StatusBadRequest)
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "400: Unable to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fmt.Printf("Accepted file: %+v, (size: %+v, mime: %+v)\n", handler.Filename, handler.Size, handler.Header)
		virus, err := scanner(file)
		if err != nil {
			fmt.Printf("scanning error reported, found %+v\n", virus)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(fmt.Sprintf("Found virus %+v\n", virus)))
		if err != nil {
			http.Error(w, "500: Error writing version response", http.StatusInternalServerError)
		}
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("version requested")
	_, err := w.Write([]byte("Version 1.1\n"))
	if err != nil {
		http.Error(w, "500: Error writing version response", http.StatusInternalServerError)
	}
}

func RunServer(scanner scanFunc) error {
	fmt.Println("Starting server listening on port :8080")
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/scan", scanHandlerWrapper(scanner)).Methods(http.MethodPost)
	api.HandleFunc("/version", versionHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	return http.ListenAndServe(":8080", nil)
}
