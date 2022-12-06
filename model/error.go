package model

type ApiError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type CustomError struct {
	Error      error `json:"error"`
	StatusCode int   `json:"statusCode"`
}
