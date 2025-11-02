package logger

type Logger interface {
	WithProperty(property Property) Logger
	WithProperties(properties []Property) Logger
	WithError(err error) Logger

	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Fatal(msg string)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Property struct {
	Key            Key
	Type           PropertyType
	IntegerValue   int
	Int32Value     int32
	Int64Value     int64
	StringValue    string
	InterfaceValue interface{}
}

type PropertyType int

const (
	// IntegerType represents an integer value on a log field.
	IntegerType PropertyType = iota
	// Int32Type represents an int32 value on a log field.
	Int32Type
	// Int64Type represents an int64 value on a log field.
	Int64Type
	// StringType represents a string value on a log field.
	StringType
)
