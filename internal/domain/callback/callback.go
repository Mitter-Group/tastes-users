package callback

type CallbackRequestBody struct {
	Code     string `json:"code"`
	State    string `json:"state"`
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

type CallbackResponse struct {
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
