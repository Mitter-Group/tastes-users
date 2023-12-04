package login

import (
	"context"

	"github.com/chunnior/users/internal/domain/callback"
)

// LoginRepository defines the methods that any data storage system must implement to get and store login data
type LoginRepository interface {
	Login(ctx context.Context, provider string) (LoginResponse, error)
	Callback(ctx context.Context, provider string) (callback.CallbackResponse, error)
}
