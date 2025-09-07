package types

type FailureResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FailureErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type NotificationError struct {
	Status  string `json:"status"` // всегда "error"
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
