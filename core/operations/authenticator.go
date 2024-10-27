package operations

import (
	"aTES/core/entities"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// This defines the service handling user related operations.
type Authenticator interface {
	CreateUser(name, role, email, joinedAt string) (int, error) // Also generates a token and stores it in the dedicated token repo.
	GetUser(userID int) (entities.User, error)
	UpdateUser(userID int, name, email, role, leftAt string) error
	DeleteUser(userID int) error
	ValidateToken(userID int, token string) bool
}

type mockAuthenticator struct {
	tokenFile string // Path to token repo yaml
	users     map[int]entities.User
}

func NewMockAuthenticator(tokenFile string) Authenticator {
	return &mockAuthenticator{
		tokenFile: tokenFile,
		users:     make(map[int]entities.User),
	}
}

func (a *mockAuthenticator) CreateUser(name, role, email, joinedAt string) (int, error) {
	newUser := entities.User{
		UserID:      len(a.users),
		Name:        name,
		Email:       email,
		Role:        role,
		Balance:     float64(0),
		JoinedAt:    joinedAt,
		LeftAt:      "",
		LastUpdated: time.Now().Format("YYYY-MM-DD HH:MM"),
	}
	a.users[len(a.users)] = newUser
	userJSON, err := json.Marshal(newUser)
	if err != nil {
		return newUser.UserID, fmt.Errorf("couldn't marshal the new user info: %w", err)
	}

	// TODO: edit url.
	_, err = http.Post("http://localhost/create-user:8181", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		return newUser.UserID, fmt.Errorf("couldn't post the new user information to the management server: %w", err)
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
// Note:
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

func (a *mockAuthenticator) DeleteUser(userID int) error

func (a *mockAuthenticator) ValidateToken(userID int, token string) bool
