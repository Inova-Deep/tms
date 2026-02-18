package http

// ErrorResponse represents a standardized API error
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
