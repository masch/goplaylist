package playlist_test

import (
	"os"
	"testing"

	"github.com/masch/goplaylist/internal/playlist"

	"github.com/stretchr/testify/require"
)

func TestListFilesSortAlphabetical(t *testing.T) {
	got, err := playlist.ListFiles("./testdata/example_1", []string{".ext", ".ext2"})
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

	got, err = playlist.ListFiles("./testdata/example_1", []string{".ext3"})
	require.NoError(t, err)
	require.EqualValues(t, emptyFilePath, got)

	got, err = playlist.ListFiles("./testdata/example_1", []string{})
	require.NoError(t, err)
	require.EqualValues(t, emptyFilePath, got)

	got, err = playlist.ListFiles("./testdata/example_2", []string{})
	require.EqualError(t, err, "lstat ./testdata/example_2: no such file or directory")
	require.EqualValues(t, emptyFilePath, got)
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

func TestPlaylistFunctional(t *testing.T) {
	// Ensure there is no ini configuration on the bootstrap and when the test finish
	// If the file doesn't exist, the error is ignored
	_ = os.Remove("cfg.ini")

	defer func() {
		require.NoError(t, os.Remove("cfg.ini"))
	}()

	got, err := playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"})
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_1/file_1_1.ext",
		"testdata/example_1/dir_1/file_1_2.ext",
		"testdata/example_1/dir_1/file_1_3.ext",
	}, got)

	got, err = playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"})
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_2/file_2_1.ext",
		"testdata/example_1/dir_2/file_2_2.ext",
		"testdata/example_1/dir_3/file_3_1.ext",
	}, got)

	got, err = playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"})
	require.NoError(t, err)
	require.EqualValues(t, []string{
		"testdata/example_1/dir_3/file_3_2.ext",
	}, got)

	var emptyList []string

	got, err = playlist.GetNextFilesFromPath("testdata/example_1", 3, []string{".ext"})
	require.NoError(t, err)
	require.EqualValues(t, emptyList, got)
}
