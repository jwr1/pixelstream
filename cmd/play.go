/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"
	"path"
	"pixelstream/internal"
	"strings"

	"github.com/spf13/cobra"
)

const playCmdExample = "  pixelstream play video.mp4 http://192.168.1.170\n  pixelstream play video.mp4 http://192.168.1.170 -r 10\n  pixelstream play video.pxlstrm http://192.168.1.170"

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:     "play SOURCE HOST",
	Short:   "Stream a video to your clock (will auto convert when needed)",
	Example: playCmdExample,
	Args:    cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		source := path.Clean(args[0])
		host, err := internal.GetUrlHost(args[1])
		cobra.CheckErr(err)

		var pixelstream *internal.PixelStream

		if playCmdNoCache {
			pixelstream, err = internal.GeneratePixelStream(source, playCmdFrameRate)
			cobra.CheckErr(err)
		} else if strings.HasSuffix(source, pixelstreamFileExt) {
			pixelstream, err = internal.LoadFile(source)
			cobra.CheckErr(err)
		} else if _, err := os.Stat(source + pixelstreamFileExt); !errors.Is(err, os.ErrNotExist) {
			pixelstream, err = internal.LoadFile(source + pixelstreamFileExt)
			cobra.CheckErr(err)
		} else {
			pixelstream, err = internal.GeneratePixelStream(source, playCmdFrameRate)
			cobra.CheckErr(err)

			err = internal.SaveFile(source+pixelstreamFileExt, pixelstream)
			cobra.CheckErr(err)
		}

		println("Version:", pixelstream.Version)
		println("Frame Rate:", pixelstream.FrameRate)
		println("# of Frames:", len(pixelstream.Frames))
		println("Duration:", pixelstream.GetDuration().String())
		println()

		pixelstream.Stream(host)
	},
}

var (
	playCmdFrameRate uint8
	playCmdNoCache   bool
)

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().Uint8VarP(&playCmdFrameRate, "frame-rate", "r", 16, "specify the frame rate when converting a video")
	playCmd.Flags().BoolVarP(&playCmdNoCache, "no-cache", "c", false, "always convert the video (instead of using cache) and don't save afterwards")
}
