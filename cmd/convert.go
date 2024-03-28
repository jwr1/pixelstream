/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path"
	"pixelstream/internal"
	"strings"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert SOURCE...",
	Short: "Pre-convert a single or multiple video files to the pxlstrm format and save",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			destDir := convertCmdOutput
			if destDir != "" {
				err := os.MkdirAll(destDir, 0755)
				cobra.CheckErr(err)
			}
			for _, source := range args {
				var dest string
				if destDir != "" {
					dest = path.Join(destDir, path.Base(source+pixelstreamFileExt))
				} else {
					dest = source + pixelstreamFileExt
				}
				convertCmdGenerateAndSave(source, dest)
			}
		} else {
			source := args[0]
			dest := convertCmdOutput

			if dest == "" {
				dest = source + pixelstreamFileExt
			} else {
				err := os.MkdirAll(path.Dir(dest), 0755)
				cobra.CheckErr(err)
				if !strings.HasSuffix(dest, pixelstreamFileExt) {
					dest += pixelstreamFileExt
				}
			}

			convertCmdGenerateAndSave(source, dest)
		}
	},
}

var (
	convertCmdFrameRate uint8
	convertCmdOutput    string
)

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&convertCmdOutput, "output", "o", "", "set the output file for a single source and output directory for multiple sources")
	convertCmd.Flags().Uint8VarP(&convertCmdFrameRate, "frame-rate", "r", 16, "specify the frame rate when generating a video")
}

func convertCmdGenerateAndSave(source string, dest string) {
	pixelstream, err := internal.GeneratePixelStream(source, convertCmdFrameRate)
	cobra.CheckErr(err)

	err = internal.SaveFile(dest, pixelstream)
	cobra.CheckErr(err)
}
