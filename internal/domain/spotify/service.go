package spotify

import (
	"github.com/chunnior/user-tastes-service/internal/domain/login"
)

type SpotifyService interface {
	Login() (login.LoginResponse, error)
}
