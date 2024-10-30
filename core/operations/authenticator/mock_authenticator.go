package Authenticator

import (
	"aTES/core/entities"
	"fmt"
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

func (a *mockAuthenticator) startServer(host string, port int) error {
	http.HandleFunc("/create_user", a.createUserHandler)
	http.HandleFunc("/get_user", a.getUserHandler)
	http.HandleFunc("/update_user", a.updateUserHandler)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}

// Returns the entities.User struct for an EXISTING user.
func (a *mockAuthenticator) getUser(userID int) (entities.User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, exists := a.users[userID]
	if !exists {
		return entities.User{}, fmt.Errorf("the user doesn't exist")
	}

	return user, nil
}

// Updates an existing user.
func (a *mockAuthenticator) updateUser(updatedUser entities.User) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Validating that the user exists.
	user, exists := a.users[updatedUser.UserID]
	if !exists {
		return fmt.Errorf("User does not exist.")
	}

	// Updating the fields of the user.
	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.Role = updatedUser.Role
	user.LeftAt = updatedUser.LeftAt
	user.LastUpdated = time.Now().Format("YYYY-MM-DD HH:MM")

	// Saving the changes.
	a.users[updatedUser.UserID] = user

	return nil
}

// Sends a delete request to remove data of a user.
func (a *mockAuthenticator) deleteUser(userID int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Validating that the user exists.
	_, exists := a.users[userID]
	if !exists {
		return fmt.Errorf("User does not exist.")
	}

	// Deteling the user.
	delete(a.users, userID)

	// Removing the token and updating the token repo.
	delete(a.tokens.tokensMap, userID)
	if err := saveTokensToYaml(a.tokens); err != nil {
		return fmt.Errorf("Error updating the token repo: %w", err)
	}

	return nil
}

func (a *mockAuthenticator) validateToken(userID int, token string) bool {
	expectedToken, exists := a.tokens.tokensMap[userID]
	if !exists || expectedToken != token {
		return false
	}

	return true
}

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
