package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
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

type rootConfig struct {
	key        string
	question   string
	defaultVal string
}

func askInput(config rootConfig) {
	if config.defaultVal == "" {
		fmt.Printf("%s\n", config.question)
	} else {
		fmt.Printf("%s [%s]\n", config.question, config.defaultVal)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		value := strings.TrimSpace(line)
		if value == "" && config.defaultVal != "" {
			viper.Set(config.key, config.defaultVal)
		} else if value != "" {
			viper.Set(config.key, value)
		}
	}
}
