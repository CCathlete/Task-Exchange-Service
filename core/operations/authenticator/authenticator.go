package Authenticator

import (
	"aTES/core/entities"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func loadTokensFromYaml(tokenYamlPath string) (tokenYaml, error) {
	var tokens tokenYaml
	tokens.location = tokenYamlPath

	data, err := os.ReadFile(tokenYamlPath)
	if err != nil {
		return tokenYaml{tokensMap: make(map[int]string)},
			fmt.Errorf("Error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &tokens)
	if err != nil {
		return tokenYaml{tokensMap: make(map[int]string)},
			fmt.Errorf("Error while loading tokens from yaml: %w", err)
	}

	return tokens, nil
}

func saveTokensToYaml(tokens tokenYaml) error {
	file, err := os.Create(tokens.location)
	if err != nil {
		return fmt.Errorf("Error recreating the file %s: %w", tokens.location, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(tokens.tokensMap); err != nil {
		return fmt.Errorf("Error writing data to fole %s: %w", tokens.location, err)
	}

	return nil
}

// TODO: Implement the GenerateToken method.
// Creates a token, refreshes the local storage of it and writes it to the token yaml file.
func (ty tokenYaml) generateToken(userID int) error

// Creates a mock authenticator from a pre declared instance and starts the server.
func InitAuthServer(host, tokenYamlPath string, port int) error {

	// Invoking the constructor and starting the server.
	var maP *mockAuthenticator
	maP, err := newMockAuthenticator(tokenYamlPath)
	if err != nil {
		return fmt.Errorf("Error starting the authenticator using yaml file at %s: %w",
			tokenYamlPath, err)
	}

	err = maP.startServer(host, port)
	if err != nil {
		return fmt.Errorf("Error starting the authentication server: %w", err)
	}

	return nil
}

// General interface functions.

func GetUser(au Authenticator, userID int) (entities.User, error)
func UpdateUser(au Authenticator, userID int, name, email, role, leftAt string) error
func DeleteUser(au Authenticator, userID int) error
func ValidateToken(au Authenticator, userID int, token string) bool
