package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/chunnior/user-tastes-service/internal/domain/callback"
	"github.com/chunnior/user-tastes-service/internal/domain/user"
	"github.com/chunnior/user-tastes-service/pkg/config"
)

type LoginRequest struct {
	Provider    string `json:"provider"`
	CallbackURL string `json:"callback_url"`
}

type LoginResponse struct {
	Url      string `json:"url"`
	Provider string `json:"provider"`
}

type LoginService struct {
	httpClient *http.Client
	config     config.Config
}

func NewLoginService(config *config.Config, httpClient *http.Client) *LoginService {
	return &LoginService{
		httpClient: httpClient,
		config:     *config,
	}
}

// Login method implementation
func (s *LoginService) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	if request.Provider == "" {
		return LoginResponse{}, errors.New("provider is required")
	}
	var loginResponse LoginResponse
	var err error
	switch request.Provider {
	case "spotify":
		loginResponse, err = callProviderLogin(s.config.SpotifyServiceURL, request.CallbackURL)
	default:
		return LoginResponse{}, errors.New("invalid provider")
	}
	if err != nil {
		return LoginResponse{}, err
	}
	loginResponse.Provider = request.Provider
	return loginResponse, nil
}

func (s *LoginService) Callback(ctx context.Context, request callback.CallbackRequestBody) (user.GenericUser, error) {
	provider := request.Provider
	if provider == "" {
		return user.GenericUser{}, errors.New("provider is required")
	}
	var genericUser user.GenericUser
	var callbackResponse callback.CallbackResponse
	//	var err error

	switch provider {
	case "spotify":
		url, err := url.Parse(s.config.SpotifyServiceURL + "/callback")
		if err != nil {
			return user.GenericUser{}, err
		}
		q := url.Query()
		q.Add("code", request.Code)
		q.Add("state", request.State)
		url.RawQuery = q.Encode()

		callbackResponse, err = callProviderCallback(url.String())
		if err != nil {
			return user.GenericUser{}, err
		}
		//	TODO: Crear objeto usuario generico
		genericUser = user.GenericUser{
			Provider:       provider,
			ProviderUserID: callbackResponse.ID,
			UserFullname:   callbackResponse.DisplayName,
			Email:          callbackResponse.Email,
			Token:          callbackResponse.Token,
			RefreshToken:   callbackResponse.RefreshToken,
		}
	default:
		return user.GenericUser{}, errors.New("invalid provider")
	}
	//	TODO: Crear usuario en BD
	genericUser.ID = "asd"
	return genericUser, nil
}

func callProviderLogin(providerUrl string, callbackUrl string) (LoginResponse, error) {
	u, err := url.Parse(providerUrl + "/login")
	if err != nil {
		return LoginResponse{}, err
	}
	q := u.Query()
	q.Add("callback_url", callbackUrl)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return LoginResponse{}, err
	}
	defer resp.Body.Close()
	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return LoginResponse{}, err
	}
	return loginResponse, nil
}

func callProviderCallback(providerUrl string) (callback.CallbackResponse, error) {
	resp, err := http.Get(providerUrl)
	if err != nil {
		return callback.CallbackResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ", err)
		return callback.CallbackResponse{}, err
	}

	// Determinar el proveedor y deserializar en la estructura apropiada
	var user callback.CallbackResponse
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("Error ", err)
		return callback.CallbackResponse{}, err
	}
	return user, nil
}
