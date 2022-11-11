package cmd

import (
	"bufio"
	"fmt"
	"github.com/artronics/vajeh-cli/internal"
	"github.com/spf13/viper"
	"os"
	osUser "os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup global config interactively",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confFile := cmd.Flags().Lookup("output").Value.String()

		user, err := osUser.Current()
		cobra.CheckErr(err)

		defaults, err := cmd.Flags().GetBool("yes")
		cobra.CheckErr(err)

		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		prompts := []internal.PromptData{
			{Key: "project_name", Label: "What the name of the project?", DefaultValue: filepath.Base(cwd)},
			{Key: "username", Label: "What is your AWS username?", DefaultValue: user.Username},
			{Key: "aws_access_key_id", Label: "What is your AWS_ACCESS_KEY_ID?", DefaultValue: os.Getenv("AWS_ACCESS_KEY_ID")},
			{Key: "aws_secret_access_key", Label: "What is your AWS_SECRET_ACCESS_KEY?", DefaultValue: os.Getenv("AWS_SECRET_ACCESS_KEY")},
			{Key: "workspace", Label: "What is your default terraform workspace? This is your short username", DefaultValue: user.Username},
			{Key: "workdir", Label: "What is the default terraform directory relative to your project path? This is where you store terraform files.", DefaultValue: "."},
			{Key: "version_prefix", Label: "What is the version prefix?", DefaultValue: "v"},
		}

		var configs map[string]string
		if defaults {
			configs = make(map[string]string, len(prompts))
			for _, p := range prompts {
				configs[p.Key] = p.DefaultValue
			}
		} else {
			configs, err = internal.GetPromptResult(prompts)
			cobra.CheckErr(err)
		}

		for k, v := range configs {
			viper.Set(k, v)
		}

		err = viper.WriteConfigAs(confFile)
		cobra.CheckErr(err)
	},
}

func setConfigValue(key string, question string, defaultVal string, setDefault bool) {
	if setDefault {
		viper.Set(key, defaultVal)
		return
	}

	if defaultVal == "" {
		fmt.Printf("%s\n", question)
	} else {
		fmt.Printf("%s [%s]\n", question, defaultVal)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		value := strings.TrimSpace(line)
		if value == "" && defaultVal != "" {
			viper.Set(key, defaultVal)
		} else if value != "" {
			viper.Set(key, value)
		}
	}

}

func scanInput(key string, question string, defaultVal string) {
	if defaultVal == "" {
		fmt.Printf("%s\n", question)
	} else {
		fmt.Printf("%s [%s]\n", question, defaultVal)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		value := strings.TrimSpace(line)
		if value == "" && defaultVal != "" {
			viper.Set(key, defaultVal)
		} else if value != "" {
			viper.Set(key, value)
		}
	}
}

func init() {
	configCmd.AddCommand(initCmd)

	confPath, err := os.UserHomeDir()
	cobra.CheckErr(err)

	confType := "yaml"
	confName := ".vajeh"
	defaultFile := fmt.Sprintf("%s/%s.%s", confPath, confName, confType)

	initCmd.Flags().StringP("output", "o", defaultFile, "the config file path")
	initCmd.PersistentFlags().BoolP("yes", "y", false, "set default values for all questions. Use ful for pipeline")
}
