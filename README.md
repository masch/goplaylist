# Goplaylist 

[![GitHub Release](https://img.shields.io/github/v/release/masch/goplaylist)](https://github.com/masch/goplaylist/releases)
[![go.dev](https://img.shields.io/badge/go.dev-reference-blue.svg)](https://pkg.go.dev/github.com/masch/goplaylist)
[![go.mod](https://img.shields.io/github/go-mod/go-version/masch/goplaylist)](go.mod)
[![Build Status](https://img.shields.io/github/workflow/status/masch/goplaylist/build)](https://github.com/masch/goplaylist/actions?query=workflow%3Abuild+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/masch/goplaylist)](https://goreportcard.com/report/github.com/masch/goplaylist)
[![codecov](https://codecov.io/gh/masch/goplaylist/branch/master/graph/badge.svg)](https://codecov.io/gh/masch/goplaylist)
[![Github Releases Stats of goplaylist](https://img.shields.io/github/downloads/masch/goplaylist/total.svg?logo=github)](https://somsubhra.com/github-release-stats/?username=masch&repository=goplaylist)

Command line application to list files from a directory path and resume from the last file used.

## Usage

`goplaylist` list files from a directory path and resume from the last file used. On every execution tracks the last file listened to resume after it on the next execution.

```
Usage: goplaylist -path=/example_path -extension=.ext_1 -extension=.ext_2 -count=3 -sort_mode=[name|timestamp_creation]

  -path string
        Specify path to load file list
  -extension value
        Specify file filter extension. Multiple extensions are supported by adding several -extension entry
  -count int
        Specify file count to load from path
  -sort_mode string
        Specify sort ascendant mode to list the files: name or timestamp_creation are supported
```

## Install

In order to install:

```bash
curl -sSfL https://raw.githubusercontent.com/masch/goplaylist/master/install.sh | sh -s --
```

## Build

- Terminal: `make` to get help for make targets.
- Terminal: `make all` to execute a full build.

## Project template 

Project based on [Golang template](https://github.com/golang-templates/seed).

## Contributing

Simply create an issue or a pull request.
