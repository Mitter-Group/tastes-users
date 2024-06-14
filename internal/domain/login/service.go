package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chunnior/users/internal/domain"
	"github.com/chunnior/users/internal/domain/callback"
	"github.com/chunnior/users/internal/domain/user"
	"github.com/chunnior/users/internal/infrastructure/aws/sqs"
	"github.com/chunnior/users/pkg/config"
	"github.com/zmb3/spotify/v2"
	v2 "google.golang.org/api/oauth2/v2"
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
	userRepo   user.UserRepository
	logger     domain.Logger
	sqsClient  *sqs.SQS
}

func NewLoginService(userRepo user.UserRepository, sqsClient *sqs.SQS, config *config.Config, httpClient *http.Client, logger domain.Logger) *LoginService {
	return &LoginService{
		httpClient: httpClient,
		config:     *config,
		userRepo:   userRepo,
		logger:     logger,
		sqsClient:  sqsClient,
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
	case "youtube":
		loginResponse, err = callProviderLogin(s.config.YoutubeServiceURL, request.CallbackURL)
	default:
		return LoginResponse{}, errors.New("invalid provider")
	}
	if err != nil {
		return LoginResponse{}, err
	}
	loginResponse.Provider = request.Provider
	return loginResponse, nil
}

func (s *LoginService) Callback(ctx context.Context, request callback.CallbackRequestBody) (*user.UserData, error) {
	provider := request.Provider
	if provider == "" {
		return &user.UserData{}, errors.New("provider is required")
	}
	var genericUser user.GenericUser
	var requestUrl string

	switch provider {
	case "spotify":
		url, err := url.Parse(s.config.SpotifyServiceURL + "/callback")
		if err != nil {
			return &user.UserData{}, err
		}
		q := url.Query()
		q.Add("code", request.Code)
		q.Add("state", request.State)
		url.RawQuery = q.Encode()
		requestUrl = url.String()
	case "youtube":
		url, err := url.Parse(s.config.YoutubeServiceURL + "/callback")
		if err != nil {
			return &user.UserData{}, err
		}
		q := url.Query()
		q.Add("code", request.Code)
		q.Add("state", request.State)
		url.RawQuery = q.Encode()
		requestUrl = url.String()
	default:
		return &user.UserData{}, errors.New("invalid provider")
	}
	genericUser, err := callProviderCallback(requestUrl, provider)
	if err != nil {
		return &user.UserData{}, err
	}
	//	Crea/actualiza usuario en BD
	savedUser, err := s.userRepo.SaveUser(ctx, genericUser)
	if err != nil {
		return &user.UserData{}, err
	}

	// envia mensaje a la queue
	err = s.createAndSendUserMessage(savedUser.ID, genericUser)
	if err != nil {
		fmt.Println("No se pudo enviar el mensaje a la queue", err)
	}
	return savedUser, nil
}

func (s *LoginService) createAndSendUserMessage(userID string, genericUser user.GenericUser) error {
	newUserMessage := user.NewUserMessage{
		UserID:         userID,
		Email:          genericUser.Email,
		Provider:       genericUser.Provider,
		ProviderUserID: genericUser.ProviderUserID,
		CreationDate:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return s.sqsClient.SendMessage(s.config.NewUserQueueURL, newUserMessage)
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

func callProviderCallback(providerUrl string, provider string) (user.GenericUser, error) {
	resp, err := http.Get(providerUrl)
	if err != nil {
		return user.GenericUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user.GenericUser{}, errors.New("error en la respuesta del proveedor")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ", err)
		return user.GenericUser{}, err
	}
	return mapProviderResponse(provider, body)
}

func mapProviderResponse(provider string, body []byte) (user.GenericUser, error) {
	var genericUser user.GenericUser
	switch provider {
	case "spotify":
		var userInfo *spotify.PrivateUser
		if err := json.Unmarshal(body, &userInfo); err != nil {
			fmt.Println("Error ", err)
			return user.GenericUser{}, err
		}
		genericUser = user.GenericUser{
			Provider:       provider,
			ProviderUserID: userInfo.ID,
			Email:          userInfo.Email,
			UserFullname:   userInfo.DisplayName,
		}
		if len(userInfo.Images) > 0 {
			genericUser.ProfilePicture = userInfo.Images[0].URL
		} else {
			genericUser.ProfilePicture = ""
		}

	case "youtube":
		var userInfo v2.Userinfo
		if err := json.Unmarshal(body, &userInfo); err != nil {
			fmt.Println("Error ", err)
			return user.GenericUser{}, err
		}
		genericUser = user.GenericUser{
			Provider:       provider,
			ProviderUserID: userInfo.Id,
			Email:          userInfo.Email,
			UserFullname:   userInfo.Name,
			ProfilePicture: userInfo.Picture,
		}
	default:
		return user.GenericUser{}, errors.New("invalid provider")
	}
	return genericUser, nil
}
