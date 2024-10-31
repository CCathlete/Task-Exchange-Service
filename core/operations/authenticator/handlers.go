package authenticator

import (
	"aTES/core/entities"
	"encoding/json"
	"fmt"
	"net/http"
)

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
		UserID   int    `json:"user_id"`
		Password string `json:"password"`
	}

	// Decoding the request's body to get userID and password.
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Validating the password.
	passwordIsValid := a.validatePassword(reqBody.UserID, reqBody.Password)
	if !passwordIsValid {
		http.Error(w, "Unauthorised acess", http.StatusUnauthorized)
		return
	}

	// Getting the user's information.
	user, err := a.getUser(reqBody.UserID)
	if err != nil {
		http.Error(w, "Error retrieving user's information.", http.StatusNotFound)
		return
	}

	// Sending user's data as json.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding user's data", http.StatusInternalServerError)
		return
	}
}

func (a *mockAuthenticator) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Making sure that we got a get request.
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Password string        `json:"name"`
		User     entities.User `json:"user"` // User has json tags of its own.
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Validating the password.
	passwordIsValid := a.validatePassword(reqBody.User.UserID, reqBody.Password)
	if !passwordIsValid {
		http.Error(w, "Unauthorised acess", http.StatusUnauthorized)
		return
	}

	// Updating the user.
	err := a.updateUser(reqBody.User)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusNotFound)
		return
	}

	// Sending a response with a success message.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully updated"))
}

func (a *mockAuthenticator) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Making sure that we got a delete request.
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Password string `json:"name"`
		UserID   int    `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Validating the password.
	passwordIsValid := a.validatePassword(reqBody.UserID, reqBody.Password)
	if !passwordIsValid {
		http.Error(w, "Unauthorised acess", http.StatusUnauthorized)
		return
	}

	// Deleting the user.
	err := a.deleteUser(reqBody.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusNotFound)
		return
	}

	// Sending a response with success message.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully updated"))
}
