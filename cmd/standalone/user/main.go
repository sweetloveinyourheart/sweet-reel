package main

import (
	"time"

	"github.com/spf13/cobra"

	appuser "github.com/sweetloveinyourheart/sweet-reel/cmd/app/services/user"
	apputils "github.com/sweetloveinyourheart/sweet-reel/cmd/app/utils"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
)

const defaultShortDescription = "Sweet Real User Service"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/sweet-reel/cmd/standalone/user generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)

	commands = append(commands, appuser.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, apputils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
