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

// arrayFlags defines custom flags to support array flags values.
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "arrays flags representations"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	if err := run(os.Args[1:], &playlist.Playlist{}); err != nil {
		log.Fatal(err)
	}
}

func run(args []string, playlistClient playlister) error {
	fileList, err := GetNextFilesFromPath(args, playlistClient)
	if err != nil {
		return err
	}

	return writeOutput(fileList)
}

func writeOutput(fileList []string) error {
	// Write the result on the std output
	writer := bufio.NewWriter(os.Stdout)

	for _, file := range fileList {
		if _, err := writer.WriteString(file + " "); err != nil {
			return err
		}
	}

	// Flush std out
	return writer.Flush()
}

// GetNextFilesFromPath get next files list from the command line using command line flags.
func GetNextFilesFromPath(args []string, playlistClient playlister) ([]string, error) {
	// parse flags values from command line
	var extensions arrayFlags

	setFlags := flag.NewFlagSet("goplaylist", flag.ContinueOnError)
	sortModeRaw := setFlags.String("short_mode", "",
		"Specify sort ascendant mode to list the files: name or timestamp_creation are supported")
	path := setFlags.String("path", "", "Specify path to load file list")
	countFiles := setFlags.Int("count", 0, "Specify file count to load from path")
	setFlags.Var(&extensions, "extension", "Specify extensions")

	if err := setFlags.Parse(args); err != nil {
		return nil, err
	}

	if *sortModeRaw == "" {
		return nil, errSortModeIsEmpty
	}

	if *path == "" {
		return nil, errPathOriginIsEmpty
	}

	if *countFiles == 0 {
		return nil, errCountFilesIsEmpty
	}

	if extensions == nil {
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
