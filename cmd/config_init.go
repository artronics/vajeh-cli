package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	osUser "os/user"
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
		if defaults {
			viper.Set("username", user.Username)
			viper.Set("aws_access_key_id", os.Getenv("AWS_ACCESS_KEY_ID"))
			viper.Set("aws_secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY"))
			viper.Set("workspace", user.Username)
			viper.Set("workdir", "")
			viper.Set("version_prefix", "v")
		} else {
			scanInput("username", "What is your AWS username?", user.Username)
			scanInput("aws_access_key_id", "What is your AWS_ACCESS_KEY_ID?", os.Getenv("AWS_ACCESS_KEY_ID"))
			scanInput("aws_secret_access_key", "What is your AWS_SECRET_ACCESS_KEY?", os.Getenv("AWS_SECRET_ACCESS_KEY"))
			scanInput("workspace", "What is your default terraform workspace? This is your short username", user.Username)
			scanInput("workdir", "What is the default terraform directory relative to your project path? This is where you store terraform files.", ".")
			scanInput("version_prefix", "What is the version prefix?", "v")
		}

		err = viper.WriteConfigAs(confFile)
		cobra.CheckErr(err)
	},
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
		} else if defaultVal != "" {
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
