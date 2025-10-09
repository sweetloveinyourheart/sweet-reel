package main

import (
	"time"

	"github.com/spf13/cobra"

	appapigateway "github.com/sweetloveinyourheart/sweet-reel/cmd/app/services/api_gateway"
	apputils "github.com/sweetloveinyourheart/sweet-reel/cmd/app/utils"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
)

const defaultShortDescription = "Sweet Real API Gateway"

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

//go:generate go run github.com/sweetloveinyourheart/sweet-reel/cmd/standalone/api_gateway generate

func main() {
	cmdutil.ServiceRootCmd.Short = defaultShortDescription
	commands := make([]*cobra.Command, 0)

	commands = append(commands, appapigateway.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, apputils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
