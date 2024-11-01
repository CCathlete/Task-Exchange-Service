package main

import (
	auth "aTES/core/operations/authenticator"
	"fmt"
	"log"
	"net/http"
)

func main() {
	passYamlPath := "/home/ccat/Repos/Task-Exchange-Service/core/operations/authenticator/passwords.yaml"
	usersYamlPath := "/home/ccat/Repos/Task-Exchange-Service/core/operations/authenticator/users.yaml"

	err := initAuthServer("localhost", passYamlPath, usersYamlPath, 8181)
	if err != nil {
		log.Fatalf("Coulden't start the authentication server: %v", err)
	}
}

// Creates a Mock authenticator from a pre declared instance and starts the server.
func initAuthServer(host, passwordYamlPath, usersYamlPath string, port int) error {

	// Invoking the constructor and starting the server.
	var maP *auth.MockAuthenticator
	maP, err := auth.NewMockAuthenticator(passwordYamlPath, usersYamlPath)
	if err != nil {
		return fmt.Errorf("error starting the authenticator using yaml file at %s: %w",
			passwordYamlPath, err)
	}

	http.HandleFunc("/create_user", maP.CreateUserHandler)
	http.HandleFunc("/get_user", maP.GetUserHandler)
	http.HandleFunc("/update_user", maP.UpdateUserHandler)
	http.HandleFunc("/delete_user", maP.DeleteUserHandler)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		return fmt.Errorf("error starting the authentication server: %w", err)
	}

	return nil
}
