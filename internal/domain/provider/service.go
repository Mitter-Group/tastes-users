package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/chunnior/users/internal/domain"
	"github.com/chunnior/users/pkg/config"
)

type ProviderService struct {
	httpClient *http.Client
	config     config.Config
	//providerRepo ProviderRepository providerRepo ProviderRepository,
	logger domain.Logger
}

func NewProviderService(config *config.Config, httpClient *http.Client, logger domain.Logger) *ProviderService {
	return &ProviderService{
		httpClient: httpClient,
		config:     *config,
		//providerRepo: providerRepo,
		logger: logger,
	}
}

func (s *ProviderService) GetProviderInfo(ctx context.Context, request ProviderInfoRequest) (DataInfoResponse, error) {
	if request.Provider == "" {
		return DataInfoResponse{}, errors.New("provider is required")
	}
	var providerInfoResponse DataInfoResponse
	var err error
	switch request.Provider {
	case "spotify":
		providerInfoResponse, err = s.callProviderInfo(s.config.SpotifyServiceURL, request)
	default:
		return DataInfoResponse{}, errors.New("invalid provider")
	}
	return providerInfoResponse, err
}

func (s *ProviderService) callProviderInfo(providerUrl string, request ProviderInfoRequest) (DataInfoResponse, error) {
	url := fmt.Sprintf("%s/%s/%s", providerUrl, request.DataType, request.UserID)
	response, err := http.Get(url)
	if err != nil {
		s.logger.Error(err.Error())
		return DataInfoResponse{}, err
	}
	defer response.Body.Close()
	var providerInfoResponse DataInfoResponse
	err = json.NewDecoder(response.Body).Decode(&providerInfoResponse)
	if err != nil {
		s.logger.Error(err.Error())
		return DataInfoResponse{}, err
	}
	return providerInfoResponse, nil
}
