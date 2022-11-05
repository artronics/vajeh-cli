package cmd

import (
	"github.com/artronics/vajeh-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create/use terraform workspace and apply",
	Long: `It creates a terraform workspace if not one already present. It switches to that workspace
and finally runs terraform apply.`,
	Run: func(cmd *cobra.Command, args []string) {
		workdir := viper.GetString("workdir")
		wss, err := internal.GetWorkspaces(workdir)
		if err != nil {
			cobra.CheckErr(err)
		}

		activeWs := wss[0]
		desiredWs := viper.GetString("workspace")

		if activeWs != desiredWs {
			err = internal.ChangeWorkspace(workdir, wss, desiredWs)
			if err != nil {
				cobra.CheckErr(err)
			}
		}

		awsCred, err := internal.GetAwsCred()
		if err != nil {
			cobra.CheckErr(err)
		}

		//fmt.Printf("Acc: %s | Sec: %s\n", awsCred.AccessKey, awsCred.AccessSecret)
		err = internal.Apply(workdir, awsCred)
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
