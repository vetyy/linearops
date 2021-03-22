package config

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetupViper initializes spf13/viper library global options
func SetupViper() error {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.AddConfigPath("/etc/linearops/") // path to look for the config file in
	viper.AddConfigPath(".")               // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return err
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("LINEAROPS")
	viper.AutomaticEnv()
	return nil
}

// SetupLogrus initializes sirupsen/logrus library global options
func SetupLogrus(logLevel string) error {
	log.SetFormatter(&log.JSONFormatter{})
	l, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	log.SetLevel(l)
	return nil
}
