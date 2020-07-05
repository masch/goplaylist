package playlist_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/masch/goplaylist/internal/playlist"
)

func TestListFilesByAlphabeticalAscSort(t *testing.T) {
	got, err := playlist.ListFilesByFileNamePath("./testdata/example_1", []string{".ext", ".ext2"})
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_1/file_1_1.ext",
		"testdata/example_1/dir_1/file_1_2.ext",
		"testdata/example_1/dir_1/file_1_3.ext",
		"testdata/example_1/dir_2/file_2_1.ext",
		"testdata/example_1/dir_2/file_2_2.ext",
		"testdata/example_1/dir_2/file_2_3.ext2",
		"testdata/example_1/dir_3/file_3_1.ext",
		"testdata/example_1/dir_3/file_3_2.ext",
	}, got)

	var emptyFilePath []string

	got, err = playlist.ListFilesByFileNamePath("./testdata/example_1", []string{".ext3"})
	require.NoError(t, err)
	require.EqualValues(t, emptyFilePath, got)

	got, err = playlist.ListFilesByFileNamePath("./testdata/example_1", []string{})
	require.NoError(t, err)
	require.EqualValues(t, emptyFilePath, got)

	got, err = playlist.ListFilesByFileNamePath("./testdata/example_2", []string{})
	require.EqualError(t, err, "lstat ./testdata/example_2: no such file or directory")
	require.EqualValues(t, emptyFilePath, got)
}

func TestListFilesByTimestampCreationAscSort(t *testing.T) {
	testCaseName := "TestListFilesByTimestampCreationAscSort"
	exampleDirectoryPath, clearFunc := createFileTimestampCreationShortTestDataExample(t, testCaseName)

	defer clearFunc()

	got, err := playlist.ListFilesByDateCreation(exampleDirectoryPath, []string{".ext", ".ext2"})
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_file_creation_short/" + testCaseName + "/00_4.ext",
		"testdata/example_file_creation_short/" + testCaseName + "/00_2.ext",
		"testdata/example_file_creation_short/" + testCaseName + "/00_1.ext",
		"testdata/example_file_creation_short/" + testCaseName + "/00_3.ext2",
	}, got)
}

func TestGetNextFiles(t *testing.T) {
	playlistExample := []string{
		"testdata/example_1/dir_1/file_1_1.ext",
		"testdata/example_1/dir_1/file_1_2.ext",
		"testdata/example_1/dir_1/file_1_3.ext",
		"testdata/example_1/dir_2/file_2_1.ext",
		"testdata/example_1/dir_2/file_2_2.ext",
		"testdata/example_1/dir_3/file_3_1.ext",
	}

	got := playlist.GetNextFiles(playlistExample, 3, "")
	require.EqualValues(t, []string{
		"testdata/example_1/dir_1/file_1_1.ext",
		"testdata/example_1/dir_1/file_1_2.ext",
		"testdata/example_1/dir_1/file_1_3.ext",
	}, got)

	got = playlist.GetNextFiles(playlistExample, 3, "testdata/example_1/dir_1/file_1_2.ext")
	require.EqualValues(t, []string{
		"testdata/example_1/dir_1/file_1_3.ext",
		"testdata/example_1/dir_2/file_2_1.ext",
		"testdata/example_1/dir_2/file_2_2.ext",
	}, got)

	got = playlist.GetNextFiles(playlistExample, 3, "testdata/example_1/dir_2/file_2_1.ext")
	require.EqualValues(t, []string{
		"testdata/example_1/dir_2/file_2_2.ext",
		"testdata/example_1/dir_3/file_3_1.ext",
	}, got)

	got = playlist.GetNextFiles(playlistExample, 3, "testdata/example_1/dir_2/file_2_2.ext")
	require.EqualValues(t, []string{
		"testdata/example_1/dir_3/file_3_1.ext",
	}, got)

	var emptyList []string

	got = playlist.GetNextFiles(playlistExample, 3, "testdata/example_1/dir_3/file_3_1.ext")
	require.EqualValues(t, emptyList, got)

	got = playlist.GetNextFiles(playlistExample, 3, "testdata/example_1/dir_1/file_3_1.ext")
	require.EqualValues(t, emptyList, got)
}

func TestPlaylistSortByUnknownMode(t *testing.T) {
	got, err := playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"}, 100000)
	require.EqualValues(t, fmt.Errorf("%w: %d", playlist.ErrUnsupportedFileSortMode, 100000), err)
	require.Empty(t, got)
}

func TestPlaylistFunctional(t *testing.T) {
	t.Run("testPlaylistSortByFileNameAscFunctional", testPlaylistSortByFileNameAscFunctional)
	t.Run("testPlaylistSortByFileTimestampCreationAscFunctional", testPlaylistSortByFileTimestampCreationAscFunctional)
}

func testPlaylistSortByFileNameAscFunctional(t *testing.T) {
	// Ensure there is no ini configuration on the bootstrap and when the test finish
	// If the file doesn't exist, the error is ignored
	_ = os.Remove("cfg.ini")

	defer func() {
		require.NoError(t, os.Remove("cfg.ini"))
	}()

	got, err := playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"}, playlist.FileSortModeFileNameAsc)
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_1/file_1_1.ext",
		"testdata/example_1/dir_1/file_1_2.ext",
		"testdata/example_1/dir_1/file_1_3.ext",
	}, got)

	got, err = playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"}, playlist.FileSortModeFileNameAsc)
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_2/file_2_1.ext",
		"testdata/example_1/dir_2/file_2_2.ext",
		"testdata/example_1/dir_3/file_3_1.ext",
	}, got)

	got, err = playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"}, playlist.FileSortModeFileNameAsc)
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_3/file_3_2.ext",
	}, got)

	var emptyList []string

	got, err = playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"}, playlist.FileSortModeFileNameAsc)
	require.NoError(t, err)
	require.EqualValues(t, emptyList, got)
}

func testPlaylistSortByFileTimestampCreationAscFunctional(t *testing.T) {
	// Ensure there is no ini configuration on the bootstrap and when the test finish
	// If the file doesn't exist, the error is ignored
	_ = os.Remove("cfg.ini")

	defer func() {
		require.NoError(t, os.Remove("cfg.ini"))
	}()

	testCaseName := "testPlaylistSortByFileTimestampCreationAscFunctional"
	exampleDirectoryPath, clearFunc := createFileTimestampCreationShortTestDataExample(t, testCaseName)

	defer clearFunc()

	const shortMode playlist.FileSortMode = playlist.FileSortModeTimestampCreationAsc
	got, err := playlist.GetNextFilesFromPath(exampleDirectoryPath, 2, []string{".ext", ".ext2"}, shortMode)
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_file_creation_short/" + testCaseName + "/00_4.ext",
		"testdata/example_file_creation_short/" + testCaseName + "/00_2.ext",
	}, got)

	got, err = playlist.GetNextFilesFromPath(exampleDirectoryPath, 2, []string{".ext", ".ext2"}, shortMode)
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_file_creation_short/" + testCaseName + "/00_1.ext",
		"testdata/example_file_creation_short/" + testCaseName + "/00_3.ext2",
	}, got)

	var emptyList []string

	got, err = playlist.GetNextFilesFromPath(exampleDirectoryPath, 2, []string{".ext", ".ext2"}, shortMode)
	require.NoError(t, err)
	require.EqualValues(t, emptyList, got)
}

func createFileTimestampCreationShortTestDataExample(t *testing.T, testCaseName string) (string, func()) {
	exampleDirectoryPath := filepath.Join("./testdata/example_file_creation_short", testCaseName)
	// Ensure there is no  example directory path on the bootstrap and when the test finish
	// If the file doesn't exist, the error is ignored
	_ = os.RemoveAll(exampleDirectoryPath)

	clearFunc := func() {
		require.NoError(t, os.RemoveAll(exampleDirectoryPath))
	}

	require.NoError(t, os.MkdirAll(exampleDirectoryPath, os.ModePerm))
	createFile(t, exampleDirectoryPath, "00_4.ext")
	createFile(t, exampleDirectoryPath, "00_0.ex3")
	createFile(t, exampleDirectoryPath, "00_2.ext")
	createFile(t, exampleDirectoryPath, "00_1.ext")
	createFile(t, exampleDirectoryPath, "00_3.ext2")

	return exampleDirectoryPath, clearFunc
}

func createFile(t *testing.T, directory string, fileName string) {
	t.Helper()

	create, err := os.Create(filepath.Join(directory, fileName))
	require.NoError(t, err)
	require.NoError(t, create.Close())
	// Since this function is used on list files by timestamp creation asc sort,
	// we add a little duration sleep in order to ensure different timestamp creation on every file
	time.Sleep(100 * time.Millisecond)
}
