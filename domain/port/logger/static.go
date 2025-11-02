package logger

var sharedInstance Logger

func RegisterSharedLoggerInstance(loggerInstance Logger) {
	sharedInstance = loggerInstance
}

func GetLogger() Logger {
	return sharedInstance
}
