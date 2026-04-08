package shared

type HttpError struct {
	Message string         `json:"message"`
	Code    uint16         `json:"code"`
	Details map[string]any `json:"details"`
}
