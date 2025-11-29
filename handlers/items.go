package handlers

import (
	"net/http"

	"github.com/integrationninjas/go-app/models"
)

func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	items := []models.Item{
		{ID: 1, Name: "Book", Price: 10.99},
		{ID: 2, Name: "Pen", Price: 2.50},
		{ID: 3, Name: "Headphones", Price: 49.99},
	}
	encodeJSON(w, items)
}
