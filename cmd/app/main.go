package main

import (
	"time"

	"github.com/spf13/cobra"

	appvideoprocessing "github.com/sweetloveinyourheart/sweet-reel/cmd/app/services/video_processing"
	apputils "github.com/sweetloveinyourheart/sweet-reel/cmd/app/utils"
	"github.com/sweetloveinyourheart/sweet-reel/pkg/cmdutil"
)

//go:generate go run github.com/sweetloveinyourheart/sweet-reel/cmd/app generate

func init() {
	// Always use UTC for generated timestamps
	time.Local = time.UTC
}

func main() {
	commands := make([]*cobra.Command, 0)

	commands = append(commands, appvideoprocessing.Command(cmdutil.ServiceRootCmd))

	commands = append(commands, apputils.CheckCommand())

	cmdutil.InitializeService(commands...)
}
