package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const pixelstreamFileExt = ".pxlstrm"

var rootCmd = &cobra.Command{
	Use:     "pixelstream SOURCE HOST",
	Short:   "Stream videos to your awtrix clock with ease.",
	Example: playCmdExample,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
