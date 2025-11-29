package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/integrationninjas/go-app/models"
)

const randomUserAPIURL = "https://randomuser.me/api/"

func GetRandomUser(w http.ResponseWriter, r *http.Request) {
	userData, err := fetchRandomUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching user data: %v", err)
		return
	}

	encodeJSON(w, userData.Results[0]) // Encode and return the first user data
}

func fetchRandomUser() (models.UserData, error) {
	resp, err := http.Get(randomUserAPIURL)
	if err != nil {
		return models.UserData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.UserData{}, err
	}

	var userData models.UserData
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return models.UserData{}, err
	}

	return userData, nil
}
