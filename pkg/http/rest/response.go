package rest

// Response holds a standard HTTP response data
type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
