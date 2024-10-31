package authenticator

import (
	"aTES/core/entities"
	"sync"
)

// This defines the service handling user related operations.
type Authenticator interface {
	CreateUser(name, role, email, joinedAt string) (int, error) // Also generates a password and stores it in the dedicated password repo.
	GetUser(userID int) (entities.User, error)
	UpdateUser(userID int, name, email, role, leftAt string) error
	DeleteUser(userID int) error
	Validatepassword(userID int, password string) bool
}

type passwordYaml struct {
	location     string         // Path to the actual yaml file.
	passwordsMap map[int]string `yaml:"passwords"`
}

type usersYaml struct {
	location string                // Path to the actual yaml file.
	usersMap map[int]entities.User `yaml:"users"`
}

type mockAuthenticator struct {
	users     *usersYaml
	passwords *passwordYaml // The yaml file is loaded here for fast drawing.
	mu        sync.Mutex
}

// type passwordRepo interface {
// 	generatepassword(userID int) error // Creates a password and pushes it to the repo.
// }
