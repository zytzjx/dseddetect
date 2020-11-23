package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// PrintHandler print all port HD
func PrintHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	ss := DetectData.String()
	w.Write([]byte(ss))
}

// LabelsHandler only show label connected
func LabelsHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	ss := DetectData.IndexString()
	w.Write([]byte(ss))
}

// LabelHandler get single port
func LabelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintf(w, "id: %v\n", vars["id"])
	ss := DetectData.ItemToString(vars["id"])
	w.Write([]byte(ss))
}
