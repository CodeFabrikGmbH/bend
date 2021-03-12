package config

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Config struct {
	Path     string    `json:"path"`
	Target   string    `json:"target"`
	Response *Response `json:"response"`
}
