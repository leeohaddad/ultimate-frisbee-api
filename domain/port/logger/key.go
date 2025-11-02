package logger

type Key string

const (
	// Success is the label that identifies if an HTTP request was handled successfully or not.
	Success Key = "Success"
	// FailureReason is the label that identifies the reason why some HTTP request failed.
	FailureReason Key = "FailureReason"
	// StatusCode is the label that identifies the response status code of an HTTP request.
	StatusCode Key = "StatusCode"
	// Key is the label that identifies the unique identifier of some entity or resource.
	ID Key = "ID"
)
