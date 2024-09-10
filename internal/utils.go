package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func GetUrlHost(rawURL string) (string, error) {
	newURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}

	if newURL.OmitHost {
		return "", errors.New("missing or invalid host")
	}

	return newURL.Scheme + "://" + newURL.Host, nil
}

func SwitchMode(mode tea.Model) (tea.Model, tea.Cmd) {
	return mode, tea.Batch(
		mode.Init(),
		tea.WindowSize(),
	)
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

var OS_FS = os.DirFS(OS_FS_ROOT)

type FileLocation struct {
	System fs.FS
	Path   string
}

func (fl FileLocation) ToOSPath() string {
	if fl.System != OS_FS {
		panic("ToOSPath() was called on a non OS FS system")
	}

	return OS_FS_ROOT + filepath.FromSlash(fl.Path)
}

func FromOSPath(path string) FileLocation {
	return FileLocation{
		System: OS_FS,
		Path:   filepath.ToSlash(strings.TrimPrefix(filepath.Clean(path), OS_FS_ROOT)),
	}
}

// ReadFile reads the named file from the file system fs and returns its contents.
// A successful call returns a nil error, not [io.EOF].
// (Because ReadFile reads the whole file, the expected EOF
// from the final Read is not treated as an error to be reported.)
//
// If fs implements [ReadFileFS], ReadFile calls fs.ReadFile.
// Otherwise ReadFile calls fs.Open and uses Read and Close
// on the returned [File].
func (fl FileLocation) ReadFile() ([]byte, error) {
	return fs.ReadFile(fl.System, fl.Path)
}

// WriteFile writes data to the file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm (before umask);
// otherwise WriteFile truncates it before writing, without changing permissions.
// Since WriteFile requires multiple system calls to complete, a failure mid-operation
// can leave the file in a partially written state.
func (fl FileLocation) WriteFile(data []byte, perm fs.FileMode) error {
	return os.WriteFile(fl.ToOSPath(), data, perm)
}
