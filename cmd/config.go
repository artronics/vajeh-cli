package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Work with configuration.",
	Long: `Manipulate default parameters by setting/resetting global configuration. The default config file should be
in "${HOME}/.vajeh.yaml". You can create one by running "vajeh config init"`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
