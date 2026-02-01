package handlers

type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewHTTPError(status int, message string) error {
	return &HTTPError{
		Status:  status,
		Message: message,
	}
}

func (e *HTTPError) Error() string {
	return e.Message
}