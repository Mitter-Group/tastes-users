package provider

type DataInfoResponse struct {
	UserID    string     `json:"user_id"`
	DataType  string     `json:"data_type"`
	Data      []DataInfo `json:"data"`
	Source    string     `json:"source"`
	Count     int        `json:"count"`
	CreatedAt string     `json:"CreatedAt"`
	UpdatedAt string     `json:"UpdatedAt"`
}

type DataInfo struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Artist      []ArtistInfo `json:"Artists"`
	ReleaseDate string       `json:"ReleaseDate"`
}

type ArtistInfo struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}
