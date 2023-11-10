package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigUnitSuite struct {
	suite.Suite
}

func TestConfigUnitSuite(t *testing.T) {
	suite.Run(t, &ConfigUnitSuite{})
}

func (us *ConfigUnitSuite) TestPrintConfig() {
	testCases := []struct {
		name           string
		config         Config
		expectedString string
	}{
		{
			name: "Config with just the log level passed",
			config: Config{
				LogLevel: "info",
			},
			expectedString: "LogLevel: info\nServerAddress: \nServerPort: \nShutDownTimeout: 0\nCORSAllowOrigins: []\nCORSAllowMethods: []\nCORSAllowHeaders: []\nCORSExposeHeaders: []\nCORSAllowCredentials: false\nCORSMaxAge: 0\n",
		},
		{
			name: "Config with multiple fields passed",
			config: Config{
				LogLevel:         "info",
				CORSAllowOrigins: []string{"http://127.0.0.1"},
				CORSAllowMethods: []string{"PUT"},
				CORSAllowHeaders: []string{"Origin content-type"},
			},
			expectedString: "LogLevel: info\nServerAddress: \nServerPort: \nShutDownTimeout: 0\nCORSAllowOrigins: [http://127.0.0.1]\nCORSAllowMethods: [PUT]\nCORSAllowHeaders: [Origin content-type]\nCORSExposeHeaders: []\nCORSAllowCredentials: false\nCORSMaxAge: 0\n",
		},
	}

	for i := range testCases {
		test := testCases[i]
		us.Run(test.name, func() {
			us.Equal(test.expectedString, test.config.MarshalConfig())
		})
	}
}
