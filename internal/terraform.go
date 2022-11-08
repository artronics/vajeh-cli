package internal

import (
	"fmt"
	"strings"
)

func execTerraform(wd string, args []string, envs []string, isStdout bool) (string, error) {
	bin := "terraform"
	chdir := fmt.Sprintf("-chdir=%s", wd)

	cmdArgs := []string{chdir}
	cmdArgs = append(cmdArgs, args...)

	out, err := Exec(bin, cmdArgs, envs, isStdout)

	if err != nil {
		if canHandleError(err.Error()) {
			if _, err := Exec(bin, []string{chdir, "init"}, envs, true); err != nil {
				return out, err
			}
			// Repeat failed command again.
			out, err := Exec(bin, cmdArgs, envs, isStdout)
			return out, err

		} else {
			return out, err
		}
	}

	return strings.TrimSpace(out), nil
}

// GetWorkspaces returns a list of all workspaces with the first item indicating the active one.
// Given correct initialized directory this list must have at least one item i.e "default"
func GetWorkspaces(wd string, credentials AwsCredentials) ([]string, error) {
	activeWs, err := execTerraform(wd, []string{"workspace", "show"}, credentials.ToEnvs(), false)
	if err != nil {
		return nil, err
	}

	var wss []string
	wss = append(wss, activeWs) // Add current one

	wsStr, err := execTerraform(wd, []string{"workspace", "list"}, credentials.ToEnvs(), false)
	if err != nil {
		return nil, err
	}
	wsAll := strings.Split(wsStr, "\n")

	for _, ws := range wsAll {
		ws := strings.TrimSpace(ws)
		if !(ws == "" || strings.HasPrefix(ws, "*")) {
			wss = append(wss, ws)
		}
	}

	return wss, nil
}

func ChangeWorkspace(wd string, credentials AwsCredentials, wss []string, ws string) error {
	// If it's already in wss then just switch (select) it otherwise create new one. "new" will switch as well
	fmt.Println(wss)
	for _, w := range wss {
		if ws == w {
			_, err := execTerraform(wd, []string{"workspace", "select", ws}, credentials.ToEnvs(), false)
			if err != nil {
				return err
			}
			return nil
		}
	}

	_, err := execTerraform(wd, []string{"workspace", "new", ws}, credentials.ToEnvs(), false)
	if err != nil {
		return err
	}

	return nil
}

func Apply(wd string, credentials AwsCredentials, vars map[string]string, isDryrun bool) error {
	var args []string
	if isDryrun {
		args = append(args, "plan")
	} else {
		args = append(args, "apply")
		args = append(args, "-auto-approve")
	}
	if varsArg := makeVarsArgs(vars); varsArg != nil {
		args = append(args, varsArg...)
	}

	_, err := execTerraform(wd, args, credentials.ToEnvs(), true)
	if err != nil {
		return err
	}

	return nil
}

func Destroy(wd string, credentials AwsCredentials, vars map[string]string, isDryrun bool) error {
	var args []string
	if isDryrun {
		args = append(args, "plan")
		args = append(args, "-destroy")
	} else {
		args = append(args, "destroy")
		args = append(args, "-auto-approve")
	}
	if varsArg := makeVarsArgs(vars); varsArg != nil {
		args = append(args, varsArg...)
	}

	_, err := execTerraform(wd, args, credentials.ToEnvs(), true)
	if err != nil {
		return err
	}

	return nil
}

func makeVarsArgs(vars map[string]string) []string {
	if len(vars) == 0 {
		return nil
	}

	var args []string
	for k, v := range vars {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}

	return args
}

// canHandleError
// These are scenarios in which "terraform init" can resolve them
func canHandleError(msg string) bool {
	return strings.Contains(msg, "please run \"terraform init\"") ||
		strings.Contains(msg, "Required plugins are not installed")
}
