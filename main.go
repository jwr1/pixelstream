package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"pixelstream/internal"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed samples
var samplesFS embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: host expected but not given. Use the following format:")
		fmt.Println("\tpixelstream <host>")
		fmt.Println("\tpixelstream http://192.168.1.170")
		os.Exit(1)
	}

	host, err := internal.GetUrlHost(os.Args[1])
	if err != nil {
		panic(err)
	}

	internal.Host = host

	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	homeDirFL := internal.FromOSPath(homeDirPath)

	samplesSubFS, err := fs.Sub(samplesFS, "samples")
	if err != nil {
		panic(err)
	}

	menuItems := []internal.MenuItem{
		{Label: "View Screen", Mode: internal.NewViewMode()},
		{Label: "Play Video", Mode: internal.NewOpenFileMode(homeDirFL.System, homeDirFL.Path)},
		{Label: "Play Sample", Mode: internal.NewOpenFileMode(samplesSubFS, ".")},
	}

	internal.MenuItems = menuItems

	_, err = tea.NewProgram(internal.NewMenuMode()).Run()
	if err != nil {
		panic(err)
	}
}
