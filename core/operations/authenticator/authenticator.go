package Authenticator

import (
	"aTES/core/entities"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
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
	tokens map[int]string `yaml:"tokens"`
}

type mockAuthenticator struct {
	users  map[int]entities.User
	tokens tokenYaml // The yaml file is loaded here for fast drawing.
	mu     sync.Mutex
}

func NewMockAuthenticator(tokenYamlPath string) (*mockAuthenticator, error) {
	tokens, err := loadTokensFromYaml(tokenYamlPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load tokens from yaml: %w", err)
	}

	return &mockAuthenticator{
		users:  make(map[int]entities.User),
		tokens: tokens,
	}, nil
}

func loadTokensFromYaml(tokenYamlPath string) (tokenYaml, error) {
	var tokens tokenYaml

	data, err := os.ReadFile(tokenYamlPath)
	if err != nil {
		return tokenYaml{tokens: make(map[int]string)},
			fmt.Errorf("Error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &tokens)
	if err != nil {
		return tokenYaml{tokens: make(map[int]string)},
			fmt.Errorf("Error while loading tokens from yaml: %w", err)
	}

	return tokens, nil
}
