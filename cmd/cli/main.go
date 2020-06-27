// Application which list files from a folder resuming the last file listed.
package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/masch/goplaylist/internal/playlist"
)

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
	// Load file list to process
	fileList, err := LoadFiles()
	if err != nil {
		log.Fatal(err)
	}

	if err := writeOutput(fileList); err != nil {
		log.Fatal(err)
	}
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
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

var (
	errPathOriginIsEmpty        = errors.New("path origin is empty")
	errCountFilesIsEmpty        = errors.New("count files is empty")
	errFilterExtensionsAreEmpty = errors.New("filter extensions are empty")
)

// LoadFiles load files list from the command line using command line flags.
func LoadFiles() ([]string, error) {
	// parse flags values from command line
	var (
		path       string
		countFiles int
		extensions arrayFlags
	)

	flag.StringVar(&path, "path", "", "Specify path to load file list")
	flag.IntVar(&countFiles, "count", 0, "Specify file count to load from path")
	flag.Var(&extensions, "extension", "Specify extensions")
	flag.Parse()

	if path == "" {
		return nil, errPathOriginIsEmpty
	}

	if countFiles == 0 {
		return nil, errCountFilesIsEmpty
	}

	if extensions == nil {
		return nil, errFilterExtensionsAreEmpty
	}

	fileList, err := playlist.GetNextFilesFromPath(path, countFiles, extensions)
	if err != nil {
		return nil, err
	}

	return fileList, nil
}
