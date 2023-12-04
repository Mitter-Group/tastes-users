package spotify

import (
	"github.com/chunnior/users/internal/domain/login"
)

type SpotifyService interface {
	Login() (login.LoginResponse, error)
}
