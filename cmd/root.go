package cmd

import (
	"fmt"
	"os"

	"github.com/itsemre/go-api-k8s/pkg/config"
	"github.com/itsemre/go-api-k8s/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/itsemre/go-api-k8s/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	Config *config.Config
	Logger *logrus.Logger
	quiet  bool
)

const (
	envType        = "env"
	envPrefix      = "API"
	configFileName = "api"
	logFormat      = "json"
)

// rootCmd is the root Cobra command of go-api-k8s
var rootCmd = &cobra.Command{
	Use:          "api",
	SilenceUsage: true,
	// Initialize everything before running child command
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Initialize the configuration
		if err = InitConfig(cmd); err != nil {
			return err
		}

		// Initialize the custom logger
		Logger, err = log.InitLogger(Config.LogLevel, logFormat)
		if err != nil {
			return err
		}
		return nil
	},
	// Run child command
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// serveCmd is the child Cobra command of go-api-k8s that begins the server
var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Begins the API",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new server instance and start it
		server := server.NewServer(Config, Logger)
		return server.Start()
	},
}

// InitConfig sets up the Viper instance and binds the Cobra flags to it
func InitConfig(cmd *cobra.Command) error {
	// Get a new data store and configuration instance
	Config = config.NewConfig()
	quiet = os.Getenv("GIN_MODE") == "release"

	// Get Viper instance
	v, err := config.GetViper(configFileName, envPrefix, envType, true, true, quiet)
	if err != nil {
		return err
	}

	// Handle the binding of parameters with one another
	if err := config.BindCobraFlags(cmd, v, envPrefix); err != nil {
		return err
	}

	// Populate the Config object with values
	if err := v.Unmarshal(&Config, viper.DecodeHook(config.DecodeSliceHook())); err != nil {
		return err
	}

	// Print configuration parameters in staging
	if !quiet {
		fmt.Printf("Configuration used:\n%s\n", Config.MarshalConfig())
	}
	return nil
}

// Execute adds child commands to the root command and sets flags appropriately.
func Execute() error {
	// Bind the created flags into the command
	if err := config.BindCobraFlagsToCmd(rootCmd, config.Config{}); err != nil {
		return err
	}
	rootCmd.AddCommand(serveCmd)
	// Execute command
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
