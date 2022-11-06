package cmd

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/artronics/vajeh-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"regexp"
)

type bumpOperation int

const (
	_ bumpOperation = iota
	major
	minor
	patch
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "manage version of a project",
	Long:  `Manage version of a project by bumping the version based of Semver`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		verPrefix := viper.GetString("version_prefix")
		withPrefix, err := cmd.Flags().GetBool("with-prefix")
		cobra.CheckErr(err)

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

		maj, err := cmd.Flags().GetBool("major")
		cobra.CheckErr(err)
		if maj {
			err = bumpVersion(releaseData, releaseFile, major, withPrefix)
			cobra.CheckErr(err)
			return
		}

		m, err := cmd.Flags().GetBool("minor")
		cobra.CheckErr(err)
		if m {
			err = bumpVersion(releaseData, releaseFile, minor, withPrefix)
			cobra.CheckErr(err)
			return
		}

		p, err := cmd.Flags().GetBool("patch")
		cobra.CheckErr(err)
		if p {
			err = bumpVersion(releaseData, releaseFile, patch, withPrefix)
			cobra.CheckErr(err)
			return
		}

		// the logic with parse-message: if it's empty or option not provided then "just print current version"
		// if message is provided but can't be parsed then, don't produce any output and exit successfully
		message, err := cmd.Flags().GetString("parse-message")
		cobra.CheckErr(err)
		if message == "" {
			printVersion(releaseData, withPrefix)
			return
		}
		op, err := parseMessage(message)
		if err != nil {
			return // Do not produce any output
		}
		err = bumpVersion(releaseData, releaseFile, op, withPrefix)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	viper.SetDefault("version_prefix", "v")

	versionCmd.PersistentFlags().Bool("major", false, "bump major version")
	versionCmd.PersistentFlags().Bool("minor", false, "bump minor version")
	versionCmd.PersistentFlags().Bool("patch", false, "bump patch version")
	versionCmd.PersistentFlags().String("parse-message", "", "parse a message (usually commit message) to decide the bump operation. If parsing fails then there wont' be any output. The pattern is `message.startsWith([release:+<bump>])` for example \"[release:+major]\"")
	versionCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "parse-message")

	versionCmd.PersistentFlags().Bool("with-prefix", false, "print the version on stdout with the prefix; useful for git tags")
}

func printVersion(rd internal.ReleaseData, withPrefix bool) {
	if withPrefix {
		fmt.Printf("%s%s\n", rd.Prefix, rd.Version)
	} else {
		fmt.Println(rd.Version)
	}
}

func bumpVersion(releaseData internal.ReleaseData, releaseFile string, op bumpOperation, withPrefix bool) error {
	version, err := semver.NewVersion(releaseData.Version)
	if err != nil {
		return err
	}

	var newVer = *version
	switch op {
	case major:
		newVer = version.IncMajor()
	case minor:
		newVer = version.IncMinor()
	case patch:
		newVer = version.IncPatch()
	default:
		return fmt.Errorf("NotSupported")
	}

	releaseData.Version = newVer.String()
	if err = releaseData.Write(releaseFile); err != nil {
		return err
	}

	printVersion(releaseData, withPrefix)

	return nil
}

const messageRegex = `^\[release:\+(major|minor|patch)\]`

func parseMessage(msg string) (bumpOperation, error) {
	r := regexp.MustCompile(messageRegex)
	s := r.FindStringSubmatch(msg)
	if s == nil {
		return -1, fmt.Errorf("ParseError")
	}
	op := s[1]
	switch op {
	case "major":
		return major, nil
	case "minor":
		return minor, nil
	case "patch":
		return patch, nil
	default:
		return -1, fmt.Errorf("ParseError")
	}
}
