package internal

import (
	"fmt"
	"strings"
)

func execTerraform(wd string, args []string) (string, error) {
	bin := "terraform"
	chdir := fmt.Sprintf("-chdir=%s", wd)

	n := []string{chdir}
	n = append(n, args...)

	out, err := Exec(bin, n)

	if err != nil {
		if strings.Contains(err.Error(), "please run \"terraform init\"") {
			_, err := Exec(bin, []string{chdir, "init"})
			if err != nil {
				return out, err
			}

		} else {
			return out, err
		}
	}

	return out, nil
}

// GetWorkspaces returns a list of all workspaces with the first item indicating the active one.
// Given correct initialized directory this list must have at least one item i.e "default"
func GetWorkspaces(wd string) ([]string, error) {
	bin := "terraform"
	chdir := fmt.Sprintf("-chdir=%s", wd)

	var wss []string
	activeWs, err := Exec(bin, []string{chdir, "workspace", "show"})
	if err != nil {
		if strings.Contains(err.Error(), "please run \"terraform init\"") {
			fmt.Printf("RUN INIT")
		}
		return nil, err
	}
	wss = append(wss, strings.TrimSpace(activeWs)) // Add current one

	wsStr, err := Exec(bin, []string{chdir, "workspace", "list"})
	if err != nil {
		if strings.Contains(err.Error(), "please run \"terraform init\"") {
			_, err := Exec(bin, []string{chdir, "init"})
			if err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
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

func ChangeWorkspace(wd string, wss []string, ws string) error {
	bin := "terraform"
	chdir := fmt.Sprintf("-chdir=%s", wd)

	// If it's already in wss then just switch (select) it otherwise create new one. "new" will switch as well
	for _, w := range wss {
		if ws == w {
			_, err := Exec(bin, []string{chdir, "workspace", "select", ws})
			if err != nil {
				return err
			}
			return nil
		}
	}

	_, err := Exec(bin, []string{"workspace", "new", ws})
	if err != nil {
		return err
	}

	return nil
}

func Apply(wd string) error {
	bin := "terraform"
	chdir := fmt.Sprintf("-chdir=%s", wd)
	s, err := Exec(bin, []string{chdir, "apply"})
	if err != nil {
		return err
	}

	fmt.Println(s)
	return nil
}
