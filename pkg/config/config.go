package config

import (
	"github.com/spf13/viper"
)

var instance *viper.Viper

func init() {
	instance = nil
}

// Instance returns the singleton instantiation of the `viper.Viper`. It should not be lazily initiated because
// it will panic if required vars are not set up correctly and needs to be run on startup.
func Instance() *viper.Viper {
	if instance == nil {
		config := viper.New()
		settings := viper.AllSettings()
		_ = config.MergeConfigMap(settings)
		instance = config
	}

	return instance
}
