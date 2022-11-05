package cmd

import (
	"errors"
	"fmt"
	"github.com/artronics/vajeh-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "manage version of a project",
	Long:  `Manage version of a project by bumping the version based of Semver`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		verPrefix := viper.GetString("version_prefix")

		releaseFile := fmt.Sprintf("%s/%s", cwd, ReleaseFile)
		if _, err := os.Stat(releaseFile); errors.Is(err, os.ErrNotExist) {
			defVer := fmt.Sprintf("%s0.1.0", verPrefix)
			fmt.Printf("%s doesn't exits. Creating one with default version %s", ReleaseFile, defVer)
			rd := internal.ReleaseData{
				Version: defVer,
			}

			err := rd.Write(releaseFile)
			cobra.CheckErr(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	viper.SetDefault("version_prefix", "v")
}
