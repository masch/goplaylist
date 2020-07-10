// Application which list files from a folder resuming the last file listed.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/masch/goplaylist/internal/playlist"
)

var (
	errSortModeIsEmpty          = errors.New("sort_mode is empty")
	errPathOriginIsEmpty        = errors.New("path origin is empty")
	errCountFilesIsEmpty        = errors.New("count files is empty")
	errFilterExtensionsAreEmpty = errors.New("filter extensions are empty")
	errUnknownFileSortMode      = errors.New("unknown file sort mode")
)

type playlister interface {
	GetNextFilesFromPath(path string, count int, fileExtension []string, sortMode playlist.FileSortMode) ([]string, error)
}

type writer interface {
	WriteString(s string) (int, error)
	Flush() error
}

// arrayFlags defines custom flags to support array flags values.
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "arrays flags representations"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var logFatal = log.Fatal //nolint // global used in order to test main result error

func main() {
	if err := run(os.Args[1:], &playlist.Playlist{}, bufio.NewWriter(os.Stdout)); err != nil {
		logFatal(err)
	}
}

func run(args []string, playlistClient playlister, playlistOutput writer) error {
	fileList, err := GetNextFilesFromPath(args, playlistClient)
	if err != nil {
		return err
	}

	return writeOutput(fileList, playlistOutput)
}

func writeOutput(fileList []string, writer writer) error {
	for _, file := range fileList {
		if _, err := writer.WriteString(file + " "); err != nil {
			return err
		}
	}

	return writer.Flush()
}

// GetNextFilesFromPath get next files list from the command line using command line flags.
// If there was an error parsing the flags arguments, it prints the usage documentation on stdout.
func GetNextFilesFromPath(args []string, playlistClient playlister) ([]string, error) {
	// parse flags values from command line
	var extensions arrayFlags

	flags := flag.NewFlagSet("goplaylist", flag.ContinueOnError)
	sortModeRaw := flags.String("sort_mode", "",
		"Specify sort ascendant mode to list the files: name or timestamp_creation are supported")
	path := flags.String("path", "", "Specify path to load file list")
	countFiles := flags.Int("count", 0, "Specify file count to load from path")
	flags.Var(&extensions, "extension",
		"Specify file filter extension. Multiple extensions are supported by adding several -extension entry")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return nil, err
	}

	if *sortModeRaw == "" {
		flags.Usage()
		return nil, errSortModeIsEmpty
	}

	if *path == "" {
		flags.Usage()
		return nil, errPathOriginIsEmpty
	}

	if *countFiles == 0 {
		flags.Usage()
		return nil, errCountFilesIsEmpty
	}

	if extensions == nil {
		flags.Usage()
		return nil, errFilterExtensionsAreEmpty
	}

	var sortMode playlist.FileSortMode

	switch *sortModeRaw {
	case "name":
		sortMode = playlist.FileSortModeFileNameAsc
	case "timestamp_creation":
		sortMode = playlist.FileSortModeTimestampCreationAsc
	default:
		return nil, fmt.Errorf("%w: %s", errUnknownFileSortMode, *sortModeRaw)
	}

	fileList, err := playlistClient.GetNextFilesFromPath(*path, *countFiles, extensions, sortMode)
	if err != nil {
		return nil, err
	}

	return fileList, nil
}
