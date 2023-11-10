package config

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AddCobraFlag creates a new flag and adds it to the flagSet
func AddCobraFlag(fs *pflag.FlagSet, ptype, name, defaultValue, help string) {
	switch ptype {
	case "int":
		intV, _ := strconv.Atoi(defaultValue)
		fs.Int(name, intV, help)
	case "bool":
		boolV, _ := strconv.ParseBool(defaultValue)
		fs.Bool(name, boolV, help)
	case "string":
		fs.String(name, defaultValue, help)
	case "[]string":
		sliceValue := strings.Split(defaultValue, " ")
		fs.StringSlice(name, sliceValue, help)
	}
}

// BindCobraFlagsToCmd creates the command-line flags based on the configuration parameters,
// and binds them to the Cobra command & the Viper instance
func BindCobraFlagsToCmd(cmd *cobra.Command, flags interface{}) error {
	// Iterate through the parameters of the struct
	for i := 0; i < reflect.TypeOf(flags).NumField(); i++ {
		parameter := reflect.TypeOf(flags).Field(i)
		name := parameter.Tag.Get("name")
		fs := cmd.Flags()
		if !cmd.HasParent() {
			fs = cmd.PersistentFlags()
		}
		// Create a flag based on the parameter and add it to the FS
		AddCobraFlag(fs,
			parameter.Type.String(),
			name,
			parameter.Tag.Get("defaultValue"),
			parameter.Tag.Get("help"))
		// Finally, bind it with Viper
		if err := viper.BindPFlag(name, fs.Lookup(name)); err != nil {
			return err
		}
	}
	return nil
}
