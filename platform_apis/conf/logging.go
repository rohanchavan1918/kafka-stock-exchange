package conf

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

// LoggingConfig specifies all the parameters needed for logging
type LoggingConfig struct {
	Level      string `viper:"string" validate:"required" mapstructure:"level"`
	File       string `viper:"string" validate:"required" mapstructure:"file"`
	Rotate     int64  `viper:"string" validate:"required" mapstructure:"rotate"`
	MaxSize    int64  `viper:"string" validate:"required" mapstructure:"max_age"`
	MaxBackups int64  `viper:"string" validate:"required" mapstructure:"max_backups"`
	MaxAge     int64  `viper:"string" validate:"required" mapstructure:"max_age"`
	Compress   bool   `viper:"string" validate:"required" mapstructure:"compress"`
}

// ConfigureLogging will take the logging configuration and also adds
// a few default parameters
func ConfigureLogging(config *LoggingConfig) (*log.Entry, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// use a file if you want
	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_CREATE|os.O_APPEND, 0660)
		if errOpen != nil {
			return nil, errOpen
		}
		log.SetOutput(bufio.NewWriter(f))
	}

	lumberjackLogrotate := &lj.Logger{
		Filename:   config.File,
		MaxSize:    int(config.MaxSize),    // Max megabytes before log is rotated
		MaxBackups: int(config.MaxBackups), // Max number of old log files to keep
		MaxAge:     int(config.MaxAge),     // Max number of days to retain log files
		Compress:   true,
	}

	if config.Level != "" {
		level, err := log.ParseLevel(strings.ToUpper(config.Level))
		if err != nil {
			return nil, err
		}
		log.SetLevel(level)
	}

	// always use the fulltimestamp
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
	})
	logMultiWriter := io.MultiWriter(os.Stdout, lumberjackLogrotate)
	log.SetOutput(logMultiWriter)

	log.WithFields(log.Fields{
		"Runtime Version": runtime.Version(),
		"Number of CPUs":  runtime.NumCPU(),
		"Arch":            runtime.GOARCH,
	}).Info("Application Initializing")

	return log.StandardLogger().WithField("hostname", hostname), nil
}
