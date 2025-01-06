package response

type APIResponse struct {
	Message string `json:"message"`
}

type APIError struct {
	Error Error `json:"error"`
}

type Error struct {
	StatusCode uint   `json:"status_code"`
	Message    string `json:"message"`
	Body       any    `json:"body,omitempty"`
}

func NewAPIError(statusCode uint, message string, body any) APIError {
	return APIError{
		Error: Error{
			StatusCode: statusCode,
			Message:    message,
			Body:       body,
		},
	}
}
