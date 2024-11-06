package authenticator

import (
	"aTES/core/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (a *MockAuthenticator) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *MockAuthenticator) GetUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *MockAuthenticator) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *MockAuthenticator) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

// Checking if the token is provided in the header, validating it and triggering login if there's no valid token.
func (a *MockAuthenticator) checkTokenAndLogin(w http.ResponseWriter, r *http.Request) (int, string, error) {
	// Extract the token fro the authorisation header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// No token was found, triggering login.
		return a.login(w, r)
	}

	// A token was found, checking if the format is Bearer <token body>.
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return 0, "", fmt.Errorf("invalid token format, expected 'Bearer <token>'")
	}

	// Validate the JWT token.
	userID, role, err := a.ValidateJWT(tokenParts[1])
	if err != nil {
		return 0, "", fmt.Errorf("invalid or expired token: %w", err)
	}

	// Token is valid, return the userID and role.
	return userID, role, nil

}

// Password authentication and generation of a new token.
func (a *MockAuthenticator) login(w http.ResponseWriter, r *http.Request) (int, string, error) {
	// Login information should contain the user credential: { "user_id": <userID>, "password": <password> }
	var loginData struct {
		UserID   int    `json:"user_id"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}

	// Parsing the login data from the request's body.
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return 0, "", fmt.Errorf("failed to decode login request body: %w", err)
	}

	// Validating the password.
	if !a.validatePassword(loginData.UserID, loginData.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return 0, "", fmt.Errorf("wrong password for user %d", loginData.UserID)
	}

	// Generatig a new JWT for the user after a successful login.
	token, err := a.GenerateJWT(loginData.UserID, loginData.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return 0, "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Responding with the token.
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))

	// Returning userID and role to be easily available for further use.
	return loginData.UserID, loginData.Role, nil
}
