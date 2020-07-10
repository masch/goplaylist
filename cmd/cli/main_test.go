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
	errProxy                 = errors.New("proxy call")
	errUndefinedFlagProvided = errors.New("flag provided but not defined: -undefined-flag")
)

type (
	input struct {
		args []string
	}

	getNextFilesFromPathReqProxy struct {
		path          string
		count         int
		fileExtension []string
		sortMode      playlist.FileSortMode
	}

	getNextFilesFromPathResProxy struct {
		fileList []string
		err      error
	}

	getNextFilesFromPathProxy struct {
		req getNextFilesFromPathReqProxy
		res getNextFilesFromPathResProxy
	}

	writerWriteStringReqProxy struct {
		string string
	}

	writerWriteStringResProxy struct {
		size int
		err  error
	}

	writerWriteString struct {
		req writerWriteStringReqProxy
		res writerWriteStringResProxy
	}

	writerFlushReqProxy struct{}

	writerFlushResProxy struct {
		err error
	}

	writerFlush struct {
		req writerFlushReqProxy
		res writerFlushResProxy
	}

	writerProxy struct {
		writeString writerWriteString
		flush       writerFlush
	}
)

func TestMainFunc(t *testing.T) {
	var fatalErrors []string

	originalLogFatal := logFatal

	defer func() {
		logFatal = originalLogFatal
	}()

	logFatal = func(v ...interface{}) {
		fatalErrors = append(fatalErrors, fmt.Sprint(v...))
	}

	main()

	for _, gotErr := range fatalErrors {
		require.Contains(t, gotErr, "flag provided but not defined: ")
	}
}

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
		suite                     suite
		getNextFilesFromPathProxy getNextFilesFromPathProxy
		writerProxies             []writerProxy
	}{
		{
			suite: suite{
				name: "FAIL_from_proxy",
				input: input{
					args: []string{"-sort_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					err: errProxy,
				},
			},
			getNextFilesFromPathProxy: getNextFilesFromPathProxy{
				req: getNextFilesFromPathReqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: getNextFilesFromPathResProxy{
					fileList: nil,
					err:      errProxy,
				},
			},
		},
		{
			suite: suite{
				name: "FAIL_from_writer_writer_string_proxy",
				input: input{
					args: []string{"-sort_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					err: errProxy,
				},
			},
			getNextFilesFromPathProxy: getNextFilesFromPathProxy{
				req: getNextFilesFromPathReqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: getNextFilesFromPathResProxy{
					fileList: []string{"file_1"},
					err:      nil,
				},
			},
			writerProxies: []writerProxy{
				{
					writeString: writerWriteString{
						req: writerWriteStringReqProxy{
							string: "file_1 ",
						},
						res: writerWriteStringResProxy{
							size: 0,
							err:  errProxy,
						},
					},
					flush: writerFlush{},
				},
			},
		},
		{
			suite: suite{
				name: "OK_with_file_timestamp_creation_sort_mode",
				input: input{
					args: []string{"-sort_mode", "timestamp_creation", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					err: nil,
				},
			},
			getNextFilesFromPathProxy: getNextFilesFromPathProxy{
				req: getNextFilesFromPathReqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeTimestampCreationAsc,
				},
				res: getNextFilesFromPathResProxy{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
			writerProxies: []writerProxy{
				{
					writeString: writerWriteString{
						req: writerWriteStringReqProxy{
							string: "file_1 ",
						},
						res: writerWriteStringResProxy{
							size: 0,
							err:  nil,
						},
					},
					flush: writerFlush{
						req: writerFlushReqProxy{},
						res: writerFlushResProxy{
							err: nil,
						},
					},
				},
				{
					writeString: writerWriteString{
						req: writerWriteStringReqProxy{
							string: "file_2 ",
						},
						res: writerWriteStringResProxy{
							size: 0,
							err:  nil,
						},
					},
					flush: writerFlush{
						req: writerFlushReqProxy{},
						res: writerFlushResProxy{
							err: nil,
						},
					},
				},
				{
					writeString: writerWriteString{
						req: writerWriteStringReqProxy{
							string: "file_3 ",
						},
						res: writerWriteStringResProxy{
							size: 0,
							err:  nil,
						},
					},
					flush: writerFlush{
						req: writerFlushReqProxy{},
						res: writerFlushResProxy{
							err: nil,
						},
					},
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
				tc.getNextFilesFromPathProxy.req.path, tc.getNextFilesFromPathProxy.req.count,
				tc.getNextFilesFromPathProxy.req.fileExtension, tc.getNextFilesFromPathProxy.req.sortMode).
				Return(tc.getNextFilesFromPathProxy.res.fileList, tc.getNextFilesFromPathProxy.res.err)

			writerMock := writerMock{}
			writerMock.Test(t)

			for _, writerProxy := range tc.writerProxies {
				writerMock.On("WriteString",
					writerProxy.writeString.req.string).
					Return(writerProxy.writeString.res.size, writerProxy.writeString.res.err)

				writerMock.On("Flush").
					Return(writerProxy.flush.res.err)
			}

			err := run(tc.suite.input.args, &playlisterMock, &writerMock)

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
		proxy getNextFilesFromPathProxy
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
			proxy: getNextFilesFromPathProxy{},
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
			proxy: getNextFilesFromPathProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_path_origin_argument",
				input: input{
					args: []string{"-sort_mode", "1"},
				},
				expect: expect{
					fileList: nil,
					err:      errPathOriginIsEmpty,
				},
			},
			proxy: getNextFilesFromPathProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_count_file_argument",
				input: input{
					args: []string{"-sort_mode", "1", "-path", "2"},
				},
				expect: expect{
					fileList: nil,
					err:      errCountFilesIsEmpty,
				},
			},
			proxy: getNextFilesFromPathProxy{},
		},
		{
			suite: suite{
				name: "FAIL_Without_filter_extension_argument",
				input: input{
					args: []string{"-sort_mode", "1", "-path", "2", "-count", "1"},
				},
				expect: expect{
					fileList: nil,
					err:      errFilterExtensionsAreEmpty,
				},
			},
			proxy: getNextFilesFromPathProxy{},
		},
		{
			suite: suite{
				name: "FAIL_With_unknown_sort_mode_argument",
				input: input{
					args: []string{"-sort_mode", "1", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: nil,
					err:      fmt.Errorf("%w: %s", errUnknownFileSortMode, "1"),
				},
			},
			proxy: getNextFilesFromPathProxy{},
		},
		{
			suite: suite{
				name: "FAIL_With_error_from_playlist_proxy",
				input: input{
					args: []string{"-sort_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: nil,
					err:      errProxy,
				},
			},
			proxy: getNextFilesFromPathProxy{
				req: getNextFilesFromPathReqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: getNextFilesFromPathResProxy{
					fileList: nil,
					err:      errProxy,
				},
			},
		},
		{
			suite: suite{
				name: "OK_with_file_name_sort_mode",
				input: input{
					args: []string{"-sort_mode", "name", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
			proxy: getNextFilesFromPathProxy{
				req: getNextFilesFromPathReqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeFileNameAsc,
				},
				res: getNextFilesFromPathResProxy{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
		},
		{
			suite: suite{
				name: "OK_with_file_timestamp_creation_sort_mode",
				input: input{
					args: []string{"-sort_mode", "timestamp_creation", "-path", "2", "-count", "1", "-extension", ".ext"},
				},
				expect: expect{
					fileList: []string{"file_1", "file_2", "file_3"},
					err:      nil,
				},
			},
			proxy: getNextFilesFromPathProxy{
				req: getNextFilesFromPathReqProxy{
					path:          "2",
					count:         1,
					fileExtension: []string{".ext"},
					sortMode:      playlist.FileSortModeTimestampCreationAsc,
				},
				res: getNextFilesFromPathResProxy{
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

type writerMock struct {
	mock.Mock
}

func (m *writerMock) WriteString(s string) (int, error) {
	args := m.Called(s)
	return args.Get(0).(int), args.Error(1)
}

func (m *writerMock) Flush() error {
	args := m.Called()
	return args.Error(0)
}

type playlisterMock struct {
	mock.Mock
}

func (m *playlisterMock) GetNextFilesFromPath(
	path string, count int, fileExtension []string, sortMode playlist.FileSortMode) ([]string, error) {
	args := m.Called(path, count, fileExtension, sortMode)
	return args.Get(0).([]string), args.Error(1)
}
