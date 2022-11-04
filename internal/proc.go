package internal

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func Exec(bin string, args []string) (string, error) {
	binary, err := exec.LookPath(bin)
	if err != nil {
		return "", fmt.Errorf("couldn't find %s executable. Make sure %s is installed", bin, bin)
	}

	cmd := exec.Command(binary, args...)

	stdout, _ := cmd.StdoutPipe()
	stderr := new(strings.Builder)
	output := new(strings.Builder)
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		output.WriteString(m)
		fmt.Println(m)
	}

	if err := cmd.Wait(); err != nil {
		if ext, ok := err.(*exec.ExitError); ok {
			// TODO: this should work for both all plat. Test it for win
			if _, ok := ext.Sys().(syscall.WaitStatus); ok {
				return "", fmt.Errorf(stderr.String())
			}
		} else {
			return "", err
		}
	}

	//stderr := new(strings.Builder)
	//stdOut := new(strings.Builder)
	//cmd.Stderr = stderr
	//cmd.Stdout = stdOut
	//
	//if err := cmd.Start(); err != nil {
	//	return "", err
	//}
	//
	//if err := cmd.Wait(); err != nil {
	//	if ext, ok := err.(*exec.ExitError); ok {
	//		// TODO: this should work for both all plat. Test it for win
	//		if _, ok := ext.Sys().(syscall.WaitStatus); ok {
	//			return "", fmt.Errorf(stderr.String())
	//		}
	//	} else {
	//		return "", err
	//	}
	//}

	return output.String(), nil
}
