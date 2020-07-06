package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/masch/goplaylist/internal/playlist"
)

var (
	errProxy                 = errors.New("proxy error")
	errUndefinedFlagProvided = errors.New("flag provided but not defined: -undefined-flag")
)

type (
	input struct {
		args []string
	}

	reqProxy struct {
		path          string
		count         int
		fileExtension []string
		sortMode      playlist.FileSortMode
	}

	resProxy struct {
		fileList []string
		err      error
	}

	dependencyProxy struct {
		req reqProxy
		res resProxy
	}
)

func TestRun(t *testing.T) { //nolint // function tool large because of BDD mechanism
	type (
		expect struct {
			err error
		}

		suite struct {
			name   string
			input  input
			expect expect
		}
	)

	tt := []struct {
		suite suite
		proxy dependencyProxy
	}{
		{
			suite: suite{
				name: "FAIL_from_proxy",
				input: input{
					args: []string{"-short_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					err: errProxy,
				},
			},
			proxy: dependencyProxy{
				req: reqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: resProxy{
					fileList: nil,
					err:      errProxy,
				},
			},
		},
		{
			suite: suite{
				name: "OK_with_file_timestamp_creation_sort_mode",
				input: input{
					args: []string{"-short_mode", "timestamp_creation", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					err: nil,
				},
			},
			proxy: dependencyProxy{
				req: reqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeTimestampCreationAsc,
				},
				res: resProxy{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.suite.name, func(t *testing.T) {
			playlisterMock := playlisterMock{}
			playlisterMock.Test(t)

			playlisterMock.On("GetNextFilesFromPath",
				tc.proxy.req.path, tc.proxy.req.count, tc.proxy.req.fileExtension, tc.proxy.req.sortMode).
				Return(tc.proxy.res.fileList, tc.proxy.res.err)

			err := run(tc.suite.input.args, &playlisterMock)

			require.EqualValues(t, tc.suite.expect.err, err)
		})
	}
}

func TestGetNextFilesFromPath(t *testing.T) { //nolint // function tool large because of BDD mechanism
	type (
		expect struct {
			fileList []string
			err      error
		}

		suite struct {
			name   string
			input  input
			expect expect
		}
	)

	tt := []struct {
		suite suite
		proxy dependencyProxy
	}{
		{
			suite: suite{
				name: "FAIL_Without_sort_mode_argument",
				input: input{
					args: []string{"--undefined-flag"},
				},
				expect: expect{
					fileList: nil,
					err:      errUndefinedFlagProvided,
				},
			},
			proxy: dependencyProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_sort_mode_argument",
				input: input{
					args: []string{},
				},
				expect: expect{
					fileList: nil,
					err:      errSortModeIsEmpty,
				},
			},
			proxy: dependencyProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_path_origin_argument",
				input: input{
					args: []string{"-short_mode", "1"},
				},
				expect: expect{
					fileList: nil,
					err:      errPathOriginIsEmpty,
				},
			},
			proxy: dependencyProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_count_file_argument",
				input: input{
					args: []string{"-short_mode", "1", "-path", "2"},
				},
				expect: expect{
					fileList: nil,
					err:      errCountFilesIsEmpty,
				},
			},
			proxy: dependencyProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_filter_extension_argument",
				input: input{
					args: []string{"-short_mode", "1", "-path", "2", "-count", "1"},
				},
				expect: expect{
					fileList: nil,
					err:      errFilterExtensionsAreEmpty,
				},
			},
			proxy: dependencyProxy{},
		},
		{
			suite: suite{
				name: "FAIL_With_unknown_short_mode_argument",
				input: input{
					args: []string{"-short_mode", "1", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: nil,
					err:      fmt.Errorf("%w: %s", errUnknownFileSortMode, "1"),
				},
			},
			proxy: dependencyProxy{},
		},
		{
			suite: suite{
				name: "FAIL_With_error_from_playlist_proxy",
				input: input{
					args: []string{"-short_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: nil,
					err:      errProxy,
				},
			},
			proxy: dependencyProxy{
				req: reqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: resProxy{
					fileList: nil,
					err:      errProxy,
				},
			},
		},
		{
			suite: suite{
				name: "OK_with_file_name_sort_mode",
				input: input{
					args: []string{"-short_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
			proxy: dependencyProxy{
				req: reqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: resProxy{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
		},
		{
			suite: suite{
				name: "OK_with_file_timestamp_creation_sort_mode",
				input: input{
					args: []string{"-short_mode", "timestamp_creation", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
			proxy: dependencyProxy{
				req: reqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeTimestampCreationAsc,
				},
				res: resProxy{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.suite.name, func(t *testing.T) {
			playlisterMock := playlisterMock{}
			playlisterMock.Test(t)

			playlisterMock.On("GetNextFilesFromPath",
				tc.proxy.req.path, tc.proxy.req.count, tc.proxy.req.fileExtension, tc.proxy.req.sortMode).
				Return(tc.proxy.res.fileList, tc.proxy.res.err)

			got, err := GetNextFilesFromPath(tc.suite.input.args, &playlisterMock)

			require.EqualValues(t, tc.suite.expect.err, err)
			require.EqualValues(t, tc.suite.expect.fileList, got)
		})
	}
}

type playlisterMock struct {
	mock.Mock
}

func (m *playlisterMock) GetNextFilesFromPath(
	path string, count int, fileExtension []string, sortMode playlist.FileSortMode) ([]string, error) {
	args := m.Called(path, count, fileExtension, sortMode)
	return args.Get(0).([]string), args.Error(1)
}
