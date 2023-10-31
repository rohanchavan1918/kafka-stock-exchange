package conf

import (
	"database/sql"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Port        int64         `viper:"int"`
	ServiceName string        `viper:"string" mapstructure:"service_name"`
	DB          DB            `mapstructure:"db"`
	Redis       Redis         `mapstructure:"redis"`
	Fluent      Fluent        `mapstructure:"fluent"`
	LogConfig   LoggingConfig `mapstructure:"log_config"`
	SlackUrl    string        `mapstructure:"slack_url"`
	KafkaConfig KafkaConfig   `mapstructure:"kafka"`
}

type appConnections struct {
	Logger      *logrus.Entry
	KafkaWriter *kafka.Writer
	DB          *sql.DB
}

var AppConnections appConnections

var AppConfig Config

func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	config := Config{}
	if err != nil {
		return nil, err
	}

	// viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile("./config/config.json")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	AppConfig = config

	return &config, nil
}
