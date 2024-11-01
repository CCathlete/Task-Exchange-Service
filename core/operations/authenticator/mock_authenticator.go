package authenticator

import (
	"aTES/core/entities"
	"fmt"
	"time"
)

func NewMockAuthenticator(passwordYamlPath, usersYamlPath string) (*MockAuthenticator, error) {
	passwords, err := loadPasswordsFromYaml(passwordYamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load passwords from yaml: %w", err)
	}
	users, err := loadUsersFromYaml(usersYamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load users from yaml: %w", err)
	}

	return &MockAuthenticator{
		users:     users,
		passwords: passwords,
	}, nil
}

// Creating a new user using the Mock authenticator.
func (a *MockAuthenticator) createUser(name, role, email, joinedAt string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	newUser := entities.User{
		UserID:      len(a.users.usersMap) + 1,
		Name:        name,
		Email:       email,
		Role:        role,
		Balance:     float64(0),
		JoinedAt:    joinedAt,
		LeftAt:      "",
		LastUpdated: time.Now().String(),
	}
	a.users.usersMap[len(a.users.usersMap)+1] = newUser

	// Updating the users yaml.
	a.users.saveUsersToYaml()

	// Generating and storing a new password for the user.
	_, err := a.newPassword(newUser.UserID)
	if err != nil {
		return newUser.UserID, fmt.Errorf("failed to create a password for user %s, %s: %w", name, role, err)
	}

	return newUser.UserID, nil
}

// Returns the entities.User struct for an EXISTING user.
func (a *MockAuthenticator) getUser(userID int) (entities.User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, exists := a.users.usersMap[userID]
	if !exists {
		return entities.User{}, fmt.Errorf("the user doesn't exist")
	}

	return user, nil
}

// Updates an existing user.
func (a *MockAuthenticator) updateUser(updatedUser entities.User) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Validating that the user exists.
	user, exists := a.users.usersMap[updatedUser.UserID]
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	// Updating the fields of the user.
	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.Role = updatedUser.Role
	user.LeftAt = updatedUser.LeftAt
	user.LastUpdated = time.Now().String()

	// Saving the changes.
	a.users.usersMap[updatedUser.UserID] = user

	return nil
}

// Sends a delete request to remove data of a user.
func (a *MockAuthenticator) deleteUser(userID int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Validating that the user exists.
	_, exists := a.users.usersMap[userID]
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	// Deleting the user and updating the user's repo.
	delete(a.users.usersMap, userID)
	if err := a.users.saveUsersToYaml(); err != nil {
		return fmt.Errorf("error updating the users repo: %w", err)
	}

	// Removing the password and updating the password repo.
	delete(a.passwords.passwordsMap, userID)
	if err := a.passwords.savePasswordsToYaml(); err != nil {
		return fmt.Errorf("error updating the password repo: %w", err)
	}

	return nil
}

func (a *MockAuthenticator) validatePassword(userID int, password string) bool {
	expectedpassword, exists := a.passwords.passwordsMap[userID]
	if !exists || expectedpassword != password {
		return false
	}

	return true
}

// Creates a new password and writes it to the password repo.
// Wrapper for passwordYaml.generatepassword
func (a *MockAuthenticator) newPassword(userID int) (string, error) {

	err := a.passwords.generatePasswordForYaml(userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate and store a new password: %w", err)
	}

	return a.passwords.passwordsMap[userID], nil
}
