package user

import "context"

type GenericUser struct {
	ID             string `json:"ID" dynamodbav:"ID"`
	Provider       string `json:"provider" dynamodbav:"Provider"`
	ProviderUserID string `json:"provider_user_id" dynamodbav:"ProviderUserID"`
	UserFullname   string `json:"user_fullname" dynamodbav:"UserFullname"`
	Email          string `json:"email" dynamodbav:"Email"`
	AccessToken    string `json:"access_token" dynamodbav:"AccessToken"`
	RefreshToken   string `json:"refresh_token" dynamodbav:"RefreshToken"`
}

type UserRepository interface {
	SaveUser(ctx context.Context, user GenericUser) (string, error)
}
