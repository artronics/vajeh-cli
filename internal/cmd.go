package internal

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetOption(cmd *cobra.Command, key string) string {
	//FIXME: use viper.BindPFlag https://github.com/spf13/viper/issues/699
	arg := cmd.Flags().Lookup(key)
	var argStr = arg.Value.String()

	confArg := viper.Get(key)
	if confArg != nil {
		argStr = confArg.(string)
	}

	// Explicitly setting argument in command line has the highest priority (overrides all comes before)
	if arg.Changed {
		argStr = arg.Value.String()
	}

	return argStr
}
