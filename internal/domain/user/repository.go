package user

import "context"

type GenericUser struct {
	// ID             string `json:"ID" dynamodbav:"ID"`
	Provider       string `json:"provider" dynamodbav:"Provider"`
	ProviderUserID string `json:"provider_user_id" dynamodbav:"ProviderUserID"`
	UserFullname   string `json:"user_fullname" dynamodbav:"UserFullname"`
	Email          string `json:"email" dynamodbav:"Email"`
}

type UserData struct {
	ID        string         `json:"ID" dynamodbav:"ID"`
	Email     string         `json:"email" dynamodbav:"Email"`
	Providers []ProviderData `json:"providers" dynamodbav:"Providers"`
}

type ProviderData struct {
	Provider       string `json:"provider" dynamodbav:"Provider"`
	ProviderUserID string `json:"provider_user_id" dynamodbav:"ProviderUserID"`
	UserFullname   string `json:"user_fullname" dynamodbav:"UserFullname"`
	Email          string `json:"email" dynamodbav:"Email"`
}

type NewUserMessage struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
	CreationDate   string `json:"creation_date"`
}

type UserRepository interface {
	SaveUser(ctx context.Context, user GenericUser) (UserData, error)
}
