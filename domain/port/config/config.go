package config

type Config interface {
	GetConfigs() (*Application, error)
}

type Application struct {
	API      *APISection
	Database *DatabaseSection
	Logger   *LoggerSection
}

type APISection struct {
	Host string
	Port int
}

type DatabaseSection struct {
	ConnectionString string
	Host             string
	Port             string
	Database         string
	User             string
	Password         string
}

type LoggerSection struct {
	Level string
	Mode  string
}
