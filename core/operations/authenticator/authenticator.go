package authenticator

import (
	"aTES/core/entities"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func loadPasswordsFromYaml(passwordYamlPath string) (*passwordYaml, error) {
	var passwords passwordYaml
	passwords.location = passwordYamlPath

	data, err := os.ReadFile(passwordYamlPath)
	if err != nil {
		return &passwordYaml{passwordsMap: make(map[int]string)},
			fmt.Errorf("error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &passwords.passwordsMap)
	if err != nil {
		return &passwordYaml{passwordsMap: make(map[int]string)},
			fmt.Errorf("error while loading passwords from yaml: %w", err)
	}

	return &passwords, nil
}

func (passwords *passwordYaml) savePasswordsToYaml() error {
	file, err := os.Create(passwords.location)
	if err != nil {
		return fmt.Errorf("error recreating the file %s: %w", passwords.location, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(passwords.passwordsMap); err != nil {
		return fmt.Errorf("error writing data to fole %s: %w", passwords.location, err)
	}

	return nil
}

func loadUsersFromYaml(usersYamlPath string) (*usersYaml, error) {
	var users usersYaml
	users.location = usersYamlPath

	data, err := os.ReadFile(usersYamlPath)
	if err != nil {
		return &usersYaml{usersMap: make(map[int]entities.User)},
			fmt.Errorf("error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &users.usersMap)
	if err != nil {
		return &usersYaml{usersMap: make(map[int]entities.User)},
			fmt.Errorf("error while loading passwords from yaml: %w", err)
	}

	return &users, nil
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

// Creates a password, refreshes the local storage of it and writes it to the password yaml file.
func (ty *passwordYaml) generatePasswordForYaml(userID int) error {

	// Generating a unique password.
	passwordBytes := make([]byte, 16) // 128 bit long password.
	if _, err := rand.Read(passwordBytes); err != nil {
		return fmt.Errorf("error generating bytes for a new password: %w", err)
	}
	newpassword := hex.EncodeToString(passwordBytes) // Converting from binary to hexadecimal and turns into a string.

	// Storing the new password in the password map fur the given userID.
	ty.passwordsMap[userID] = newpassword

	// Updating the password repo (yaml).

	if err := ty.savePasswordsToYaml(); err != nil {
		return fmt.Errorf("error while undating password repo with the new password: %w", err)
	}

	return nil
}

// General interface functions.

// func GetUser(au Authenticator, userID int) (entities.User, error)
// func UpdateUser(au Authenticator, userID int, name, email, role, leftAt string) error
// func DeleteUser(au Authenticator, userID int) error
// func Validatepassword(au Authenticator, userID int, password string) bool
