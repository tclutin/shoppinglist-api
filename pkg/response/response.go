package response

type APIError[T any] struct {
	Error Error[T] `json:"error"`
}

type Error[T any] struct {
	StatusCode uint   `json:"status_code"`
	Message    string `json:"message"`
	Body       *T     `json:"body,omitempty"`
}

func NewAPIError[T any](statusCode uint, message string, body *T) APIError[T] {
	return APIError[T]{
		Error: Error[T]{
			StatusCode: statusCode,
			Message:    message,
			Body:       body,
		},
	}
}
