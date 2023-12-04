package spotify

type SpotifyCallbackRequestProperties struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func (s SpotifyCallbackRequestProperties) GetCode() string {
	return s.Code
}

func (s SpotifyCallbackRequestProperties) GetState() string {
	return s.State
}

type SpotifyUser struct {
	DisplayName  string       `json:"display_name"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Followers    Followers    `json:"followers"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Images       []Image      `json:"images"`
	URI          string       `json:"uri"`
	Country      string       `json:"country"`
	Email        string       `json:"email"`
	Product      string       `json:"product"`
	Birthdate    string       `json:"birthdate"`
}

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

type Followers struct {
	Total int    `json:"total"`
	Href  string `json:"href"`
}

type Image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

func (su SpotifyUser) GetID() string {
	return su.ID
}

// GetEmail implementa UserResponse para SpotifyUser.
func (su SpotifyUser) GetEmail() string {
	return su.Email
}
