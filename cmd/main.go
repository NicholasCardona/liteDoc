package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"jsonCache/pkg/db"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	// load data previously saved
	err := db.LoadData("./data.gob")
	if err != nil {
		log.Println("No DB file found")
	}

	r := mux.NewRouter()

	r.HandleFunc("/store", handlePost).Methods("POST")
	r.HandleFunc("/store/{id}", handleGet).Methods("GET")
	r.HandleFunc("/store/{id}", handleDelete).Methods("DELETE")

	fmt.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	db.Store[id] = data

	db.SaveData("./data.gob")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, found := db.Store[id]
	if !found {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(data)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, exists := db.Store[id]; !exists {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	delete(db.Store, id)
	db.SaveData("./sdata.gob")

	w.WriteHeader(http.StatusNoContent)
}
