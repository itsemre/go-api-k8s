package config

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// GetViper sets up the Viper boilerplate configuration. i.e. reads the config from configFileName, sets the environment
// variables' prefix to envPrefix and the environment file type to envType; returns a pointer to (already set up) Viper.
func GetViper(configFileName, envPrefix, envType string, allowEmpty, autoEnv, quiet bool) (*viper.Viper, error) {
	v := viper.New()
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	configDir := path.Join(home, configDirName)
	v.AddConfigPath(configDir)
	v.SetConfigName(configFileName)
	v.SetConfigType(envType)
	v.SetEnvPrefix(envPrefix)
	v.AllowEmptyEnv(allowEmpty)
	if autoEnv {
		v.AutomaticEnv()
	}
	v.SetEnvKeyReplacer((strings.NewReplacer("-", "_")))

	if err = v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		if !quiet {
			fmt.Printf("Configuration file %s in path %s not found\n", configFileName, configDir)
		}
	}
	return v, nil
}

// DecodeSliceHook performs the parsing of a []string flag during the unmarshaling process of Viper
func DecodeSliceHook() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type, // data type
		t reflect.Type, // target data type
		data interface{}, // raw data
	) (interface{}, error) {
		// If the data is of type []string parse it correctly and return
		if t == reflect.TypeOf([]string{}) {
			str := data.(string)
			str = strings.TrimPrefix(str, "[")
			str = strings.TrimSuffix(str, "]")
			return strings.Split(str, ","), nil
		}
		return data, nil
	}
}
