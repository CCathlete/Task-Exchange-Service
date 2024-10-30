package Authenticator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: add handler functions for createUser, updateUser, deleteUser, getUser and an init funtion that will invoke the constructor.

func (a *mockAuthenticator) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Making sure that we get a post request.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		Email    string `json:"email"`
		JoinedAt string `json:"joined_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Using the fields of the temporary struct in the createUser method.
	userID, err := a.createUser(reqBody.Name, reqBody.Role, reqBody.Email, reqBody.JoinedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	// Sending a response with the new user's ID.
	response := map[string]int{"user_id": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *mockAuthenticator) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Making sure that we got a get request.
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		UserID int    `json:"user_id"`
		Token  string `json:"token"`
	}

	// Decoding the request's body to get userID and token.
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

}
