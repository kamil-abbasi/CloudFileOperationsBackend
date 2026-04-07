package shared

type HttpError struct {
	Message string
	Code    uint16
	Details map[string]any
}
