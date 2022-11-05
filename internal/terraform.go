package internal

import (
	"fmt"
	"strings"
)

func execTerraform(wd string, args []string, envs []string) (string, error) {
	bin := "terraform"
	chdir := fmt.Sprintf("-chdir=%s", wd)

	n := []string{chdir}
	n = append(n, args...)

	out, err := Exec(bin, n, envs)

	if err != nil {
		if strings.Contains(err.Error(), "please run \"terraform init\"") {
			_, err := Exec(bin, []string{chdir, "init"}, envs)
			if err != nil {
				return out, err
			}

		} else {
			return out, err
		}
	}

	return strings.TrimSpace(out), nil
}

// GetWorkspaces returns a list of all workspaces with the first item indicating the active one.
// Given correct initialized directory this list must have at least one item i.e "default"
func GetWorkspaces(wd string) ([]string, error) {
	activeWs, err := execTerraform(wd, []string{"workspace", "show"}, nil)
	if err != nil {
		return nil, err
	}

	var wss []string
	wss = append(wss, activeWs) // Add current one

	wsStr, err := execTerraform(wd, []string{"workspace", "list"}, nil)
	wsAll := strings.Split(wsStr, "\n")

	for _, ws := range wsAll {
		ws := strings.TrimSpace(ws)
		if !(ws == "" || strings.HasPrefix(ws, "*")) {
			wss = append(wss, ws)
		}
	}

	return wss, nil
}

func ChangeWorkspace(wd string, wss []string, ws string) error {
	// If it's already in wss then just switch (select) it otherwise create new one. "new" will switch as well
	for _, w := range wss {
		if ws == w {
			_, err := execTerraform(wd, []string{"workspace", "select", ws}, nil)
			if err != nil {
				return err
			}
			return nil
		}
	}

	_, err := execTerraform(wd, []string{"workspace", "new", ws}, nil)
	if err != nil {
		return err
	}

	return nil
}

func Apply(wd string, credentials AwsCredentials) error {
	_, err := execTerraform(wd, []string{"apply", "-auto-approve"}, credentials.ToEnvs())
	if err != nil {
		return err
	}

	return nil
}
