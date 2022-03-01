package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func GetApp(fSys afero.Fs) *cli.App {
	a := &cli.App{
		Name:        "mbox-tools",
		Description: "Tools for working with mbox files",
		Version:     fmt.Sprintf("%v, commit %v, built at %v", version, commit, date),
	}

	a.Commands = []*cli.Command{
		newSplitCommand(fSys),
	}
	return a
}

func withMBox(fSys afero.Fs, fileName string, cb func(io.Reader) error) error {
	file, err := fSys.Open(fileName)
	if err != nil {
		return err
	}
	err = cb(file)
	file.Close()
	return err
}
