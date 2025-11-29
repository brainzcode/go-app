package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/integrationninjas/go-app/models"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	message := models.Message{Text: "Hello from GoLang API!"}
	encodeJSON(w, message)
}

func encodeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error encoding data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
