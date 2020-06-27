// Package playlist list files from a folder resuming the last file listed.
package playlist

import (
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/ini.v1"
)

const (
	_iniFileName                     = "cfg.ini"
	_iniLastFileNameProcessedSection = "last"
)

// GetNextFilesFromPath returns existing files names on the path given following the next steps:
// 1. List file names and filter them by the extension given.
// 2. Load from the ini configuration file which was the last file name processed.
// If there is no file, it will return empty string.
// 3. Get next N count value given file from the last file name processed.
// 4. Save the last file name returned on the filter list.
// 5. Return the full list to processed.
func GetNextFilesFromPath(path string, count int, fileExtension []string) ([]string, error) {
	// List files from the path given
	fileList, err := ListFiles(path, fileExtension)
	if err != nil {
		return nil, err
	}

	// If there is not files, return empty list
	if len(fileList) == 0 {
		return nil, nil
	}

	// Load ini file
	cfg, err := ini.LooseLoad(_iniFileName)
	if err != nil {
		return nil, err
	}
	// Tries to load the last file name processed
	lastFileNameUsed := cfg.Section(path).Key(_iniLastFileNameProcessedSection).String()

	// If the last file name used if the same of the last file list, it means that there is no more file to list
	if lastFileNameUsed == fileList[len(fileList)-1] {
		return nil, nil
	}

	// Get n count file names after the last file name used
	nextFiles := GetNextFiles(fileList, count, lastFileNameUsed)
	if len(nextFiles) == 0 {
		return nil, nil
	}

	// Save the last file used on the ini configuration
	lastFile := nextFiles[len(nextFiles)-1]
	cfg.Section(path).Key(_iniLastFileNameProcessedSection).SetValue(lastFile)

	if err := cfg.SaveTo(_iniFileName); err != nil {
		return nil, err
	}

	return nextFiles, nil
}

// ListFiles lists file path on the given path and filter them with extension given.
func ListFiles(path string, filterExtensions []string) ([]string, error) {
	var paths []string

	// Walks on the path given finding all the find names and filter them by the extension given
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore if it is a directory
		if f.IsDir() {
			return nil
		}

		// Filter file by the extension given
		for _, fileExtension := range filterExtensions {
			if filepath.Ext(f.Name()) == fileExtension {
				paths = append(paths, path)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort file path alphabetically
	sort.Slice(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})

	return paths, nil
}

// GetNextFiles return the count given file path names from the file list given after the from the file path given.
func GetNextFiles(fileList []string, count int, fromFilePath string) []string {
	var filePaths []string

	// If the from file path given is empty, the first files given are returned
	if fromFilePath == "" {
		for i := 0; i < len(fileList); i++ {
			// Append file path on the result
			filePaths = append(filePaths, fileList[i])
			// Return file path collected when the size count is equal the count given
			if len(filePaths) == count {
				return filePaths
			}
		}

		return filePaths
	}

	// Walks over the file list given
	for i := 0; i < len(fileList); i++ {
		// Find file name position of the file name given
		if fromFilePath == fileList[i] {
			// Get the next items count given
			for j := 0; count >= j; j++ {
				i++
				// Return file paths collected when there is no more files to walk
				if i >= len(fileList) {
					return filePaths
				}
				// Append file path on the result
				filePaths = append(filePaths, fileList[i])
				// Return file path collected when the size count is equal the count given
				if len(filePaths) == count {
					return filePaths
				}
			}
		}
	}

	return filePaths
}
