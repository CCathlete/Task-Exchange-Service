package authenticator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// Initialising the mock authenticator.
	auth, err := newMockAuthenticator("~/Repos/Task-Exchange-Service/core/operations/authenticator/tokens.yaml")
	if err != nil {
		t.Errorf("Error creating a new authenticator instance: %v", err)
	}

	// Preparing the request's body.
	reqBody := `{"name":"Ken Cat", "role":"admin", "email":"kctest@example.com", "joined_at":"2024-01-01"}`

	// Creating a new request.
	req := httptest.NewRequest(http.MethodPost, "/create_user", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Calling the handler.
	auth.createUserHandler(w, req)

	// Checking the status code.
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status ok, got %v", response.Status)
	}
}
