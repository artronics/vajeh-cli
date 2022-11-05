package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vajeh",
	Short: "A cli tool to deploy code into aws using terraform",
	Long: `A tool to perform aws deployment using terraform. The deployment takes advantage of 
terraform workspace to manage deployment environment`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vajeh.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deployCmd.PersistentFlags().StringP("workspace", "w", "default", "the terraform workspace")
	err := viper.BindPFlag("workspace", deployCmd.PersistentFlags().Lookup("workspace"))
	cobra.CheckErr(err)
	viper.SetDefault("workspace", "default")

	deployCmd.PersistentFlags().StringP("workdir", "d", ".", "working directory i.e. where terraform files are located")
	err = viper.BindPFlag("workdir", deployCmd.PersistentFlags().Lookup("workdir"))
	cobra.CheckErr(err)
	viper.SetDefault("workdir", ".")

	// key and secret must be provided via environment variables only
	err = viper.BindEnv("aws-secret-access-key", "aws-access-key-id")
	if err != nil {
		cobra.CheckErr(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		confPath, err := os.UserHomeDir()
		cobra.CheckErr(err)

		confType := "yaml"
		confName := ".vajeh"
		file := fmt.Sprintf("%s/%s.%s", confPath, confName, confType)
		if _, err = os.Stat(file); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Config file %s doesn't exist. Run \"vajeh config init\" to set it up.\n", file)
			// TODO: We should exit but uncommenting below line will cause error during "vajeh config init"!
			//os.Exit(1)
		}

		viper.AddConfigPath(confPath)
		viper.SetConfigType(confType)
		viper.SetConfigName(confName)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// TODO: Same as above todo. We should exit but uncommenting below line will cause error during "vajeh config init"!
		//	fmt.Printf("Error loading config file: %s\n Make sure file exists or run \"vajeh config init\" to create config file.", err)
		//os.Exit(1)
	}
}
