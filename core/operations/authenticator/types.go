package Authenticator

import (
	"aTES/core/entities"
	"sync"
)

// This defines the service handling user related operations.
type Authenticator interface {
	CreateUser(name, role, email, joinedAt string) (int, error) // Also generates a token and stores it in the dedicated token repo.
	GetUser(userID int) (entities.User, error)
	UpdateUser(userID int, name, email, role, leftAt string) error
	DeleteUser(userID int) error
	ValidateToken(userID int, token string) bool
}

type tokenYaml struct {
	location  string         // Path to the actual yaml file.
	tokensMap map[int]string `yaml:"tokens"`
}

type mockAuthenticator struct {
	users  map[int]entities.User
	tokens tokenYaml // The yaml file is loaded here for fast drawing.
	mu     sync.Mutex
}

type tokenRepo interface {
	GenerateToken(userID int) error // Creates a token and pushes it to the repo.
}
