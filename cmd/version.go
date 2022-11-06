package cmd

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
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

		// Check if release file exists and if not create a default one
		if _, err := os.Stat(releaseFile); errors.Is(err, os.ErrNotExist) {
			defVer := "0.1.0"
			fmt.Printf("%s doesn't exits. Creating one with default version %s%s", ReleaseFile, verPrefix, defVer)
			rd, err := internal.NewReleaseData(defVer, verPrefix)
			cobra.CheckErr(err)
			err = rd.Write(releaseFile)
			cobra.CheckErr(err)
		}

		releaseData, err := internal.ParseReleaseFile(releaseFile)
		cobra.CheckErr(err)
		version, err := semver.NewVersion(releaseData.Version)
		cobra.CheckErr(err)

		var newVer = *version // set new ver to current version so, we print the current version in case of no bump option

		major, err := cmd.Flags().GetBool("major")
		cobra.CheckErr(err)
		if major {
			newVer = version.IncMajor()
		}

		minor, err := cmd.Flags().GetBool("minor")
		cobra.CheckErr(err)
		if minor {
			newVer = version.IncMinor()
		}

		patch, err := cmd.Flags().GetBool("patch")
		cobra.CheckErr(err)
		if patch {
			newVer = version.IncPatch()
		}

		releaseData.Version = newVer.String()
		err = releaseData.Write(releaseFile)
		cobra.CheckErr(err)

		withPrefix, err := cmd.Flags().GetBool("with-prefix")
		cobra.CheckErr(err)

		if withPrefix {
			fmt.Printf("%s%s\n", releaseData.Prefix, releaseData.Version)
		} else {
			fmt.Println(releaseData.Version)
		}

		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	viper.SetDefault("version_prefix", "v")

	versionCmd.PersistentFlags().Bool("major", false, "bump major version")
	versionCmd.PersistentFlags().Bool("minor", false, "bump minor version")
	versionCmd.PersistentFlags().Bool("patch", false, "bump patch version")
	versionCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch")

	versionCmd.PersistentFlags().Bool("with-prefix", false, "print the version on stdout with the prefix; useful for git tags")
}
