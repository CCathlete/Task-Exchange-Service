package operations

import "aTES/core/entities"

// This defines the service handling user related operations.
type Authenticator interface {
	CreateUser(name, role, joinedAt string) (int, error) // Also generates a token and stores it in the dedicated token repo.
	GetUser(userID int) (entities.User, error)
	UpdateUser(userID int, name, email, role, leftAt string) error
	DeleteUser(userID int) error
	ValidateToken(userID int, token string) bool
}

type mockAuthenticator struct {
	tokenFile string // Path to token repo yaml
	users     map[int]entities.User
}

func NewMockAuthenticator(tokenFile string) Authenticator {
	return &mockAuthenticator{
		tokenFile: tokenFile,
		users:     make(map[int]entities.User),
	}
}

func (a *mockAuthenticator) CreateUser(name, role, joinedAt string) (int, error)

func (a *mockAuthenticator) GetUser(userID int) (entities.User, error)

func (a *mockAuthenticator) UpdateUser(userID int, name, email, role, leftAt string) error

func (a *mockAuthenticator) DeleteUser(userID int) error

func (a *mockAuthenticator) ValidateToken(userID int, token string) bool
