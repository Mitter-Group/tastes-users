package google

type GoogleCallbackRequestProperties struct {
	Code  string `json:"code"`
	State string `json:"state"`
	Token string `json:"token"`
}

func (g GoogleCallbackRequestProperties) GetCode() string {
	return g.Code
}

func (g GoogleCallbackRequestProperties) GetState() string {
	return g.State
}

type GoogleCallbackResponse struct {
	DisplayName string `json:"display_name"`
	//	Followers   Followers `json:"followers"`
	Href      string `json:"href"`
	ID        string `json:"id"`
	URI       string `json:"uri"`
	Country   string `json:"country"`
	Email     string `json:"email"`
	Product   string `json:"product"`
	Birthdate string `json:"birthdate"`
}

func (g GoogleCallbackResponse) GetDisplayName() string {
	return g.DisplayName
}
