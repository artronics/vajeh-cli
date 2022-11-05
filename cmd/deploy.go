package cmd

import (
	"fmt"
	"github.com/artronics/vajeh-cli/internal"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create/use terraform workspace and apply",
	Long: `It creates a terraform workspace if not one already present. It switches to that workspace
and finally runs terraform apply.`,
	Run: func(cmd *cobra.Command, args []string) {
		workdir := internal.GetOption(cmd, "workdir")
		wss, err := internal.GetWorkspaces(workdir)
		if err != nil {
			cobra.CheckErr(err)
		}

		activeWs := wss[0]
		desiredWs := internal.GetOption(cmd, "workspace")

		fmt.Printf("all: %s\nfinal desired: %s | active: %s\n", wss, desiredWs, activeWs)
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

		fmt.Printf("Acc: %s | Sec: %s\n", awsCred.AccessKey, awsCred.AccessSecret)
		//err = internal.Apply(workdir, awsCred)
		//if err != nil {
		//	cobra.CheckErr(err)
		//}
		fmt.Println(workdir)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("workspace", "w", "default", "the terraform workspace")
	deployCmd.Flags().StringP("workdir", "d", ".", "working directory i.e. where terraform files are located")
}
