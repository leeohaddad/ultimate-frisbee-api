package logger

type LogLevel string

const (
	// DebugLevel is the log level that shows every log message.
	DebugLevel LogLevel = "debug"
	// InfoLevel is the log level that shows every log message of level info or above.
	InfoLevel LogLevel = "info"
	// WarnLevel is the log level that shows every log message of level warn or above.
	WarnLevel LogLevel = "warn"
	// ErrorLevel is the log level that shows every log message of level error or above.
	ErrorLevel LogLevel = "error"
	// FatalLevel is the log level that shows every log message of level fatal.
	FatalLevel LogLevel = "fatal"
)
