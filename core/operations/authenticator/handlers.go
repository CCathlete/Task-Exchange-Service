package Authenticator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: add handler functions for createUser, updateUser, deleteUser, getUser and an init funtion that will invoke the constructor.

func (a *mockAuthenticator) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Making sure that we send a post request.
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

// Creates a nre mock authenticator from a pre declared instance and starts the server.
// The server's host address, port and the path to the token yaml file are inside the request's body.
func (a *mockAuthenticator) initAuthServer(w http.ResponseWriter, r *http.Request) {

	// Making sure that we send a get request.
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Host          string `json:"host"`
		Port          int    `json:"port"`
		TokenYamlPath string `json:"token_yaml_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Now we have the data we need to invoke the constructor and starting the server.
	var maP *mockAuthenticator
	maP, err := newMockAuthenticator(reqBody.TokenYamlPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating an authenticator from yaml in %s: %v",
			reqBody.TokenYamlPath, err), http.StatusInternalServerError)
		return
	}

	err = maP.startServer(reqBody.Host, reqBody.Port)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting the auth server: %v",
			err), http.StatusInternalServerError)
		return
	}
	// TODO: we need to find a place to use this as an handle function since all other handle functions are in startServer.
}
