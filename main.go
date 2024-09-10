package main

import (
	"embed"
	"io/fs"
	"os"
	"pixelstream/internal"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed samples
var samplesFS embed.FS

func main() {
	host, err := internal.GetUrlHost(os.Args[1])
	if err != nil {
		panic(err)
	}

	internal.Host = host

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	samplesSubFS, err := fs.Sub(samplesFS, "samples")
	if err != nil {
		panic(err)
	}

	menuItems := []internal.MenuItem{
		{Label: "View Screen", Mode: internal.NewViewMode()},
		{Label: "Play Video", Mode: internal.NewOpenFileMode(os.DirFS("/"), strings.TrimPrefix(homeDir, "/"))},
		{Label: "Play Sample", Mode: internal.NewOpenFileMode(samplesSubFS, ".")},
	}

	internal.MenuItems = menuItems

	_, err = tea.NewProgram(internal.NewMenuMode()).Run()
	if err != nil {
		panic(err)
	}
}
