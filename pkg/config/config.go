package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	sensitiveTag  = "sensitive"
	configDirName = ".api"
)

// Config is the object that holds all the configuration parameters of the server, and holds all the information necessary to create command-line flags for them
type Config struct {
	LogLevel             string   `mapstructure:"LOG_LEVEL" name:"log-level" long:"log-level" defaultValue:"info" help:"Logging level, can only be one of 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'."`
	ServerAddress        string   `mapstructure:"SERVER_ADDRESS" name:"server-address" long:"server-address" defaultValue:"127.0.0.1" help:"The address that the web server will be listening to"`
	ServerPort           string   `mapstructure:"SERVER_PORT" name:"server-port" long:"server-port" defaultValue:"8080" help:"The port that the web server will be listening to"`
	ShutDownTimeout      int      `mapstructure:"SHUTDOWN_TIMEOUT" name:"shutdown-timeout" long:"shutdown-timeout" defaultValue:"10" help:"The timeout (in seconds) for the server to shut down"`
	CORSAllowOrigins     []string `mapstructure:"CORS_ALLOW_ORIGINS" name:"cors-allow-origins" long:"cors-allow-origins" defaultValue:"*" help:"Allow origins for CORS configuration"`
	CORSAllowMethods     []string `mapstructure:"CORS_ALLOW_METHODS" name:"cors-allow-methods" long:"cors-allow-methods" defaultValue:"GET POST PUT DELETE" help:"List of CORS methods that are allowed"`
	CORSAllowHeaders     []string `mapstructure:"CORS_ALLOW_HEADERS" name:"cors-allow-headers" long:"cors-allow-headers" defaultValue:"Origin content-type" help:"List of CORS headers that are allowed"`
	CORSExposeHeaders    []string `mapstructure:"CORS_EXPOSE_HEADERS" name:"cors-expose-headers" long:"cors-expose-headers" defaultValue:"Content-Length" help:"List of CORS headers that are exposed"`
	CORSAllowCredentials bool     `mapstructure:"CORS_ALLOW_CREDENTIALS" name:"cors-allow-credentials" long:"cors-allow-credentials" defaultValue:"false" help:"Whether to allow credentials to CORS"`
	CORSMaxAge           int      `mapstructure:"CORS_MAX_AGE" name:"cors-max-age" long:"cors-max-age" defaultValue:"1" help:"Maximum age (in hours) pertaining to CORS configuration"`
}

// NewConfig returns an instance of the Config
func NewConfig() *Config {
	return &Config{}
}

// MarshalConfig stringifies the configuration parameters while excluding sensitive ones
func (c *Config) MarshalConfig() string {
	var sb strings.Builder
	v := reflect.ValueOf(*c)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// Extract the "sensitive" tag from the struct and skip it if it's sensitive
		structField := typeOfS.Field(i)
		f := structField.Tag.Get(sensitiveTag)
		if f == "" {
			switch structField.Type.Kind() {
			case reflect.Bool:
				w := fmt.Sprintf("%s: %t\n", typeOfS.Field(i).Name, v.Field(i).Interface())
				sb.WriteString(w)
			case reflect.Int:
				w := fmt.Sprintf("%s: %d\n", typeOfS.Field(i).Name, v.Field(i).Interface())
				sb.WriteString(w)
			default:
				w := fmt.Sprintf("%s: %s\n", typeOfS.Field(i).Name, v.Field(i).Interface())
				sb.WriteString(w)
			}
		}
	}
	return sb.String()
}

// BindCobraFlags offers the integration between Cobra and Viper, through binding of the flags defined
// in a *cobra.Command to the corresponding configuration set to *viper.Viper
func BindCobraFlags(cmd *cobra.Command, v *viper.Viper, envPrefix string) error {
	var (
		err    error
		envVar string
	)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables don't have dashes in their names, so bind them to their equivalent
		// keys with underscores, e.g. --allow-origin to ENVPREFIX_ALLOW_ORIGIN
		if strings.Contains(f.Name, "-") {
			envVar = strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		} else {
			envVar = f.Name
		}
		if bindErr := v.BindEnv(envVar); bindErr != nil {
			// fxn signature of the callback passed to VisitAll does not return an error
			// therefore, we exit abruptly here
			err = fmt.Errorf("error when binding flag %s to its corresponding env var\n%s", f.Name, bindErr)
			return
		}
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if f.Changed {
			v.Set(envVar, f.Value.String())
		} else if !f.Changed && !v.IsSet(envVar) {
			v.Set(envVar, f.DefValue)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
