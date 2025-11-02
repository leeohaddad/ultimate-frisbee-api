package zap

import (
	"fmt"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	logger *zap.Logger
	fields []zap.Field
	err    error
}

func NewLogger() (*Zap, error) {
	zapLogger, err := buildLogger()

	return &Zap{
		logger: zapLogger,
	}, err
}

func buildLogger() (*zap.Logger, error) {
	logLevel := zap.InfoLevel

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Encoding = "json"
	loggerConfig.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("error when building logger from config: %w", err)
	}

	return logger, nil
}

func (zapLogger *Zap) WithProperty(property logger.Property) logger.Logger {
	return &Zap{
		logger: zapLogger.logger,
		fields: append(zapLogger.fields, mapLoggerField(property)),
		err:    zapLogger.err,
	}
}

func (zapLogger *Zap) WithProperties(properties []logger.Property) logger.Logger {
	return &Zap{
		logger: zapLogger.logger,
		fields: buildLoggerFields([]zap.Field{}, properties, nil),
		err:    zapLogger.err,
	}
}

func (zapLogger *Zap) WithError(err error) logger.Logger {
	return &Zap{
		logger: zapLogger.logger,
		fields: zapLogger.fields,
		err:    err,
	}
}

func (zapLogger *Zap) Debugf(format string, args ...interface{}) {
	logMessage := buildMessage(fmt.Sprintf(format, args...), zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Debug(logMessage, zapFields...)
}

func (zapLogger *Zap) Debug(msg string) {
	logMessage := buildMessage(msg, zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Debug(logMessage, zapFields...)
}

func (zapLogger *Zap) Infof(format string, args ...interface{}) {
	logMessage := buildMessage(fmt.Sprintf(format, args...), zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Info(logMessage, zapFields...)
}

func (zapLogger *Zap) Info(msg string) {
	logMessage := buildMessage(msg, zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Info(logMessage, zapFields...)
}

func (zapLogger *Zap) Errorf(format string, args ...interface{}) {
	logMessage := buildMessage(fmt.Sprintf(format, args...), zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Error(logMessage, zapFields...)
}

func (zapLogger *Zap) Error(msg string) {
	logMessage := buildMessage(msg, zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Error(logMessage, zapFields...)
}

func (zapLogger *Zap) Fatalf(format string, args ...interface{}) {
	logMessage := buildMessage(fmt.Sprintf(format, args...), zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Fatal(logMessage, zapFields...)
}

func (zapLogger *Zap) Fatal(msg string) {
	logMessage := buildMessage(msg, zapLogger.err)
	zapFields := buildLoggerFields(zapLogger.fields, []logger.Property{}, zapLogger.err)
	zapLogger.logger.Fatal(logMessage, zapFields...)
}

func buildMessage(msg string, err error) string {
	if err != nil {
		return fmt.Sprintf("%s: %s", msg, err.Error())
	}

	return msg
}

func buildLoggerFields(baseFields []zap.Field, fields []logger.Property, err error) []zap.Field {
	zapFields := append([]zap.Field{}, baseFields...)

	for _, field := range fields {
		zapFields = append(zapFields, mapLoggerField(field))
	}
	if err != nil {
		zapFields = append(zapFields, zap.Field{
			Key:    "error",
			Type:   zapcore.StringType,
			String: err.Error(),
		})
	}

	return zapFields
}

func mapLoggerField(property logger.Property) zap.Field {
	return zap.Field{
		Key:       string(property.Key),
		Type:      mapLoggerFieldType(property.Type),
		Integer:   mapLoggerFieldInteger(property.Int32Value, property.Int64Value),
		String:    property.StringValue,
		Interface: property.InterfaceValue,
	}
}

func mapLoggerFieldType(fieldType logger.PropertyType) zapcore.FieldType {
	switch fieldType {
	case logger.IntegerType:
		return zapcore.Int32Type
	case logger.Int32Type:
		return zapcore.Int32Type
	case logger.Int64Type:
		return zapcore.Int64Type
	case logger.StringType:
		return zapcore.StringType
	default:
		return zapcore.UnknownType
	}
}

func mapLoggerFieldInteger(integer32 int32, integer64 int64) int64 {
	if integer32 != 0 {
		return int64(integer32)
	}

	return integer64
}
