package provider

import "context"

type ProviderInfoRequest struct {
	Provider string
	UserID   string
	DataType string
}

type ProviderRepository interface {
	GetProviderInfo(ctx context.Context, request ProviderInfoRequest) (string, error)
}
