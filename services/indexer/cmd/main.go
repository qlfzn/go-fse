package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var f File
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("received file: %+v\n", f)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"success","message":"File indexed successfully"}`)
}

func main() {
	http.HandleFunc("/index", indexHandler)
	log.Println("mock indexer running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
