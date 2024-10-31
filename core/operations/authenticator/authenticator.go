package authenticator

import (
	"aTES/core/entities"
	"crypto/rand"
	"encoding/hex"
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
			fmt.Errorf("error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &tokens.tokensMap)
	if err != nil {
		return tokenYaml{tokensMap: make(map[int]string)},
			fmt.Errorf("error while loading tokens from yaml: %w", err)
	}

	return tokens, nil
}

func (tokens *tokenYaml) saveTokensToYaml() error {
	file, err := os.Create(tokens.location)
	if err != nil {
		return fmt.Errorf("error recreating the file %s: %w", tokens.location, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(tokens.tokensMap); err != nil {
		return fmt.Errorf("error writing data to fole %s: %w", tokens.location, err)
	}

	return nil
}

func loadUsersFromYaml(usersYamlPath string) (usersYaml, error) {
	var users usersYaml
	users.location = usersYamlPath

	data, err := os.ReadFile(usersYamlPath)
	if err != nil {
		return usersYaml{usersMap: make(map[int]entities.User)},
			fmt.Errorf("error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &users.usersMap)
	if err != nil {
		return usersYaml{usersMap: make(map[int]entities.User)},
			fmt.Errorf("error while loading tokens from yaml: %w", err)
	}

	return users, nil
}

func (users *usersYaml) saveUsersToYaml() error {
	file, err := os.Create(users.location)
	if err != nil {
		return fmt.Errorf("error recreating the file %s: %w", users.location, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(users.usersMap); err != nil {
		return fmt.Errorf("error writing data to fole %s: %w", users.location, err)
	}

	return nil
}

// Creates a token, refreshes the local storage of it and writes it to the token yaml file.
func (ty *tokenYaml) generateTokenForYaml(userID int) error {

	// Generating a unique token.
	tokenBytes := make([]byte, 16) // 128 bit long token.
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("error generating bytes for a new token: %w", err)
	}
	newToken := hex.EncodeToString(tokenBytes) // Converting from binary to hexadecimal and turns into a string.

	// Storing the new tiken in the token map fur the given userID.
	ty.tokensMap[userID] = newToken

	// Updating the token repo (yaml).

	if err := ty.saveTokensToYaml(); err != nil {
		return fmt.Errorf("error while undating token repo with the new token: %w", err)
	}

	return nil
}

// Creates a mock authenticator from a pre declared instance and starts the server.
func InitAuthServer(host, tokenYamlPath, usersYamlPath string, port int) error {

	// Invoking the constructor and starting the server.
	var maP *mockAuthenticator
	maP, err := newMockAuthenticator(tokenYamlPath, usersYamlPath)
	if err != nil {
		return fmt.Errorf("error starting the authenticator using yaml file at %s: %w",
			tokenYamlPath, err)
	}

	err = maP.startServer(host, port)
	if err != nil {
		return fmt.Errorf("error starting the authentication server: %w", err)
	}

	return nil
}

// General interface functions.

// func GetUser(au Authenticator, userID int) (entities.User, error)
// func UpdateUser(au Authenticator, userID int, name, email, role, leftAt string) error
// func DeleteUser(au Authenticator, userID int) error
// func ValidateToken(au Authenticator, userID int, token string) bool
