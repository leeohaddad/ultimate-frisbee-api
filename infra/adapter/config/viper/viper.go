package viper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"
	"github.com/spf13/viper"
)

const configPrefix = "UF_API"

type Viper struct {
	configurations *config.Application
}

// NewConfig builds the config based on a YAML config file.
func NewConfig(configFile string) (*Viper, error) {
	configs, err := buildConfigurations(configFile)
	if err != nil {
		return nil, fmt.Errorf("error while building configs: %w", err)
	}

	return &Viper{
		configurations: configs,
	}, nil
}

// GetConfigs returns the configuration object.
func (viperConfig *Viper) GetConfigs() (*config.Application, error) {
	if viperConfig.configurations == nil {
		return nil, errors.New("viper configurations object is undefined")
	}

	return viperConfig.configurations, nil
}

// buildConfig builds a config structure based on a YAML config file.
func buildConfigurations(configFile string) (*config.Application, error) {
	viperConfig, err := loadConfiguration(configFile)
	if err != nil {
		return nil, err
	}

	appConfig := mapConfigurations(viperConfig)

	return appConfig, nil
}

func loadConfiguration(configFile string) (*viper.Viper, error) {
	viperConfig := viper.New()
	if configFile != "" { // enable ability to specify config file via flag
		viperConfig.SetConfigFile(configFile)
	}
	viperConfig.SetEnvPrefix(configPrefix)
	viperConfig.SetConfigType("yaml")
	viperConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperConfig.AutomaticEnv()

	if err := viperConfig.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not load configuration file: %w", err)
	}

	return viperConfig, nil
}

func mapConfigurations(viperConfig *viper.Viper) *config.Application {
	// set defaults if needed
	viperConfig.SetDefault("api.host", "0.0.0.0")
	viperConfig.SetDefault("api.port", "42000")
	viperConfig.SetDefault("zap.level", "ERROR")

	// map Viper structure into an Application config, encapsulating Viper logic here
	return &config.Application{
		API: &config.APISection{
			Host: viperConfig.GetString("api.host"),
			Port: viperConfig.GetInt("api.port"),
		},
		Database: &config.DatabaseSection{
			ConnectionString: viperConfig.GetString("database.connectionString"),
		},
		Logger: &config.LoggerSection{
			Level: viperConfig.GetString("zap.level"),
			Mode:  viperConfig.GetString("zap.mode"),
		},
	}
}
