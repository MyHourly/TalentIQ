package dto

// APIResponse is the common response structure returned by the API.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is used when an API request fails.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
