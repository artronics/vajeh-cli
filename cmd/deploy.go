package cmd

import (
	"fmt"
	"github.com/artronics/vajeh-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create/use terraform workspace and apply",
	Long: `It creates a terraform workspace if not one already present. It switches to that workspace
and finally runs terraform apply.`,
	Run: func(cmd *cobra.Command, args []string) {
		workdir := viper.GetString("workdir")
		wss, err := internal.GetWorkspaces(workdir)
		cobra.CheckErr(err)

		activeWs := wss[0]
		desiredWs := viper.GetString("workspace")

		awsCred, err := internal.GetAwsCred()
		cobra.CheckErr(err)

		if activeWs != desiredWs {
			err = internal.ChangeWorkspace(workdir, awsCred, wss, desiredWs)
			cobra.CheckErr(err)
		}

		isDryrun, err := cmd.Flags().GetBool("dryrun")
		cobra.CheckErr(err)
		isDestroy, err := cmd.Flags().GetBool("destroy")
		cobra.CheckErr(err)
		varsArg, err := cmd.Flags().GetString("vars")
		cobra.CheckErr(err)

		vars, err := parseVars(varsArg)
		cobra.CheckErr(err)

		if isDestroy {
			err = internal.Destroy(workdir, awsCred, vars, isDryrun)
		} else {
			err = internal.Apply(workdir, awsCred, vars, isDryrun)
		}
		cobra.CheckErr(err)
	},
}

func parseVars(vars string) (map[string]string, error) {
	m := make(map[string]string)
	if vars == "" {
		return m, nil
	}

	ss := strings.Fields(vars)
	for _, s := range ss {
		kv := strings.Split(s, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("syntax error while parsing vars. It must be in the for of: \"<key1>:<value1> <key2>:<value2>\"")
		}
		m[kv[0]] = kv[1]
	}

	return m, nil
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringP("workspace", "w", "default", "the terraform workspace")
	err := viper.BindPFlag("workspace", deployCmd.PersistentFlags().Lookup("workspace"))
	cobra.CheckErr(err)
	viper.SetDefault("workspace", "default")

	deployCmd.PersistentFlags().StringP("workdir", "d", ".", "working directory i.e. where terraform files are located")
	err = viper.BindPFlag("workdir", deployCmd.PersistentFlags().Lookup("workdir"))
	cobra.CheckErr(err)
	viper.SetDefault("workdir", ".")

	deployCmd.PersistentFlags().Bool("dryrun", false, "whether to apply the plan. It's equivalent of terraform plan")
	deployCmd.PersistentFlags().Bool("destroy", false, "whether to destroy deployment. It WONT ask for confirmation; add --dryrun along this option to check the plan")
	deployCmd.PersistentFlags().String("vars", "", "terraform variable in the form of: --vars \"<key1>:<value1> <key2>:<value2>\"")
}
