package Authenticator

import (
	"aTES/core/entities"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func newMockAuthenticator(tokenYamlPath string) (*mockAuthenticator, error) {
	tokens, err := loadTokensFromYaml(tokenYamlPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load tokens from yaml: %w", err)
	}

	return &mockAuthenticator{
		users:  make(map[int]entities.User),
		tokens: tokens,
	}, nil
}

// Creating a new user using the mock authenticator.
func (a *mockAuthenticator) createUser(name, role, email, joinedAt string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	newUser := entities.User{
		UserID:      len(a.users) + 1,
		Name:        name,
		Email:       email,
		Role:        role,
		Balance:     float64(0),
		JoinedAt:    joinedAt,
		LeftAt:      "",
		LastUpdated: time.Now().Format("YYYY-MM-DD HH:MM"),
	}
	a.users[len(a.users)+1] = newUser

	// Generating and storing a new token for the user.
	_, err := a.newToken(newUser.UserID)
	if err != nil {
		return newUser.UserID, fmt.Errorf("Failed to create a token for user %s, %s: %w", name, role, err)
	}

	return newUser.UserID, nil
}

func (a *mockAuthenticator) GetUser(userID int) (entities.User, error) {
	user, exists := a.users[userID]
	if !exists {
		return entities.User{}, fmt.Errorf("the user doesn't exist")
	}
	// If we reach this line it means user isn't stored locally or doesn't exist yet.

	// TODO: edit url for get request.
	response, err := http.Get("http://localhost/get:8181")
	if err != nil {
		return entities.User{}, fmt.Errorf("couldn't get user via http: %w", err)
	}
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return entities.User{}, fmt.Errorf("error while reading the response's body: %w", err)
	}

	err = json.Unmarshal(resBody, &user)
	if err != nil {
		return entities.User{}, fmt.Errorf("error while unmarshaling the response body's json: %w", err)
	}

	return user, nil
}

// Updates an existing user.
func (a *mockAuthenticator) UpdateUser(userID int, name, email, role, leftAt string) error {
	// Validating that the user exists.
	user, exists := a.users[userID]
	if !exists {
		return fmt.Errorf("User does not exist.")
	}

	// Updating the fields of the user.
	user.Name = name
	user.Email = email
	user.Role = role
	user.LeftAt = leftAt
	user.LastUpdated = time.Now().Format("YYYY-MM-DD HH:MM")
	a.users[userID] = user

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Couldnt marshal user data into a JSON: %w", err)
	}

	// Create a request body and send it via an http client.
	request, err := http.NewRequest(http.MethodPut, "http://localhost/update:8181", bytes.NewBuffer(userJSON))
	if err != nil {
		return fmt.Errorf("Error when creating a new http put request from the user's JSON: %w", err)
	}

	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("Error when sending an http put request to the user mgmt server: %w", err)
	}

	return nil
}

// Sends a delete request to remove data of a user.
func (a *mockAuthenticator) deleteUser(userID int) error

func (a *mockAuthenticator) validateToken(userID int, token string) bool

// Creates a new token and writes it to the token repo.
// Wrapper for tokenYaml.generateToken
func (a *mockAuthenticator) newToken(userID int) (string, error) {
	// TODO: add type check to check if tokenRepo is of type tokenYaml.

	err := a.tokens.generateToken(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to generate and store a new token: %w", err)
	}

	return a.tokens.tokensMap[userID], nil
}
