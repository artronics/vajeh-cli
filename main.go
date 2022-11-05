package main

import (
	"github.com/artronics/vajeh-cli/cmd"
	"github.com/artronics/vajeh-cli/internal"
)

func main() {
	_, _ = internal.Exec("ls", []string{}, []string{}, true)
	cmd.Execute()
}
