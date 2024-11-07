package authenticator

import (
	"aTES/core/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (a *MockAuthenticator) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Calling checkTokenAndLogin to validate the login credentials + token or trigger a login.
	loginUserID, loginRole, err := a.checkTokenAndLogin(w, r)

	// Checking if the logged in user exists.
	if _, exists := a.users.usersMap[loginUserID]; !exists {
		http.Error(w, fmt.Sprintf("Login does not exist: %v", err), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorised: %v", err), http.StatusUnauthorized)
		return
	}

	// Checking that the login was made by an admin.
	if loginRole != "admin" {
		http.Error(w, "Forbidden: unauthorised access.", http.StatusForbidden)
		return
	}

	// Making sure that we get a post request.
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct { // Target user information sits inside the request's body.
		Target struct {
			Name     string `json:"name"`
			Role     string `json:"role"`
			Email    string `json:"email"`
			JoinedAt string `json:"joined_at"`
		} `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Using the fields of the temporary struct in the createUser method.
	userID, err := a.createUser(reqBody.Target.Name, reqBody.Target.Role, reqBody.Target.Email, reqBody.Target.JoinedAt)
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
	// Calling checkTokenAndLogin to validate the login credentials + token or trigger a login.
	loginUserID, _, err := a.checkTokenAndLogin(w, r)

	// Checking if the logged in user exists.
	if _, exists := a.users.usersMap[loginUserID]; !exists {
		http.Error(w, fmt.Sprintf("Login does not exist: %v", err), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorised: %v", err), http.StatusUnauthorized)
		return
	}

	// Making sure that we got a get request.
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Target struct {
			UserID   int    `json:"user_id"`
			Password string `json:"password"`
		}
	}

	// Decoding the request's body to get userID and password.
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Validating the password.
	passwordIsValid := a.validatePassword(reqBody.Target.UserID, reqBody.Target.Password)
	if !passwordIsValid {
		http.Error(w, "Unauthorised acess", http.StatusUnauthorized)
		return
	}

	// Getting the user's information.
	user, err := a.getUser(reqBody.Target.UserID)
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
	// Calling checkTokenAndLogin to validate the login credentials + token or trigger a login.
	loginUserID, loginRole, err := a.checkTokenAndLogin(w, r)

	// Checking if the logged in user exists.
	if _, exists := a.users.usersMap[loginUserID]; !exists {
		http.Error(w, fmt.Sprintf("Login does not exist: %v", err), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorised: %v", err), http.StatusUnauthorized)
		return
	}

	// Checking that the login was made by an admin.
	if loginRole != "admin" {
		http.Error(w, "Forbidden: unauthorised access.", http.StatusForbidden)
		return
	}

	// Making sure that we got a get request.
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Target struct {
			Password string        `json:"name"`
			User     entities.User `json:"user"` // User has json tags of its own.
		}
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Validating the password.
	passwordIsValid := a.validatePassword(reqBody.Target.User.UserID, reqBody.Target.Password)
	if !passwordIsValid {
		http.Error(w, "Unauthorised acess", http.StatusUnauthorized)
		return
	}

	// Updating the user.
	err = a.updateUser(reqBody.Target.User)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusNotFound)
		return
	}

	// Sending a response with a success message.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully updated"))
}

func (a *MockAuthenticator) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Calling checkTokenAndLogin to validate the login credentials + token or trigger a login.
	loginUserID, loginRole, err := a.checkTokenAndLogin(w, r)

	// Checking if the logged in user exists.
	if _, exists := a.users.usersMap[loginUserID]; !exists {
		http.Error(w, fmt.Sprintf("Login does not exist: %v", err), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorised: %v", err), http.StatusUnauthorized)
		return
	}

	// Checking that the login was made by an admin.
	if loginRole != "admin" {
		http.Error(w, "Forbidden: unauthorised access.", http.StatusForbidden)
		return
	}

	// Making sure that we got a delete request.
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// Creating a temporary structure for decoding the request body's json.
	var reqBody struct {
		Target struct {
			Password string `json:"name"`
			UserID   int    `json:"user_id"`
		}
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error decoding the request's body: %w", http.StatusBadRequest)
		return
	}

	// Validating the password.
	passwordIsValid := a.validatePassword(reqBody.Target.UserID, reqBody.Target.Password)
	if !passwordIsValid {
		http.Error(w, "Unauthorised acess", http.StatusUnauthorized)
		return
	}

	// Deleting the user.
	err = a.deleteUser(reqBody.Target.UserID)
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
	// Extract the token from the authorisation header.
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
		// Triggering login if token is not valid/ expired.
		return a.login(w, r)
	}

	// Token is valid, return the userID and role.
	return userID, role, nil

}

// Password authentication and generation of a new token.
func (a *MockAuthenticator) login(w http.ResponseWriter, r *http.Request) (int, string, error) {
	// Login information should contain the user credential: { "user_id": <userID>, "password": <password> }
	var loginData struct {
		Login struct {
			UserID   int    `json:"user_id"`
			Role     string `json:"role"`
			Password string `json:"password"`
		} `json:"login"`
	}

	// Parsing the login data from the request's body.
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return 0, "", fmt.Errorf("failed to decode login request body: %w", err)
	}

	// Validating the password.
	if !a.validatePassword(loginData.Login.UserID, loginData.Login.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return 0, "", fmt.Errorf("wrong password for user %d", loginData.Login.UserID)
	}

	// Generatig a new JWT for the user after a successful login.
	token, err := a.GenerateJWT(loginData.Login.UserID, loginData.Login.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return 0, "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Responding with the token.
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))

	// Returning userID and role to be easily available for further use.
	return loginData.Login.UserID, loginData.Login.Role, nil
}
