package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/mirtchovski/clamav"
	"github.com/gorilla/mux"
)

func scanHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("scan requested")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("version requested")
	_, err := w.Write([]byte("Version 1.1\n"))
	if err != nil {
		http.Error(w, "500: Error writing version response", http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/scan", scanHandler).Methods(http.MethodPost)
	api.HandleFunc("/version", versionHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
