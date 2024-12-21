package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfigWithPath used to load config object from specified path.
// The config filename loaded is depends on provided environmentMode.
// Currently supported config file extension is "json", "toml", "yaml", "yml", "properties", "props", "prop" and "hcl".
// eg:
// - provided environmentMode is "stg", so it will load "{path}/stg.{supported_extension}" file.
// - provided environmentMode is "prod", so it will load "{path}/prod.{supported_extension}" file.
func LoadConfigWithPath(configObject interface{}, environtmentMode string, path string) error {
	viper.SetConfigName(environtmentMode)
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("load config error. Cause: %v", err)
	}

	err = viper.Unmarshal(&configObject)
	if err != nil {
		return fmt.Errorf("load config error. Cause: %v", err)
	}

	return nil
}

func LoadConfig(configObject interface{}, environtmentMode string) error {
	return LoadConfigWithPath(configObject, environtmentMode, "./config")
}
