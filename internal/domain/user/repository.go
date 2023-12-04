package user

type GenericUser struct {
	ID             string `json:"ID"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
	UserFullname   string `json:"user_fullname"`
	Email          string `json:"email"`
	Token          string `json:"token"`
	RefreshToken   string `json:"refresh_token"`
}
