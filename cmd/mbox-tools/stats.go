package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"

	"github.com/aquilax/mbox-tools/lib/mbox"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func newStatsCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:      "stats",
		Usage:     "Mbox file statistics",
		ArgsUsage: "[filename]...",
		Subcommands: []*cli.Command{
			newStatsSenderCommand(fSys),
			newStatsSenderSizeCommand(fSys),
		},
	}
}

func newStatsSenderCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:      "sender",
		Usage:     "Prints out number of messages per sender",
		ArgsUsage: "[filename]...",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "ignore-errors",
				Usage: "Overwrite target files",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			ignoreErrors := c.Bool("ignore-errors")
			stats := make(map[string]int, 0)
			for _, sourceFile := range c.Args().Slice() {
				err := withMBox(fSys, sourceFile, func(mb io.Reader) error {
					index := 0
					err := mbox.ReadMessages(mb, func(b []byte) (bool, error) {
						r := bufio.NewReader(bytes.NewBuffer(b))
						sender, err := getSender(r)
						if err != nil && !ignoreErrors {
							return true, fmt.Errorf("message # %d: %v", index, err)
						}
						if v, found := stats[sender]; found {
							stats[sender] = v + 1
						} else {
							stats[sender] = 1
						}
						index++
						return false, nil
					})
					return err
				})
				if err != nil {
					return fmt.Errorf("file: %s: %v", sourceFile, err)
				}
			}
			return printReport(stats, os.Stdout)
		},
	}
}

func newStatsSenderSizeCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:      "sender-size",
		Usage:     "Prints out number of messages per sender",
		ArgsUsage: "[filename]...",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "ignore-errors",
				Usage: "Overwrite target files",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			ignoreErrors := c.Bool("ignore-errors")
			stats := make(map[string]int, 0)
			for _, sourceFile := range c.Args().Slice() {
				err := withMBox(fSys, sourceFile, func(mb io.Reader) error {
					index := 0
					err := mbox.ReadMessages(mb, func(b []byte) (bool, error) {
						size := len(b)
						r := bufio.NewReader(bytes.NewBuffer(b))
						sender, err := getSender(r)
						if err != nil && !ignoreErrors {
							return true, fmt.Errorf("message # %d: %v", index, err)
						}
						if v, found := stats[sender]; found {
							stats[sender] = v + size
						} else {
							stats[sender] = size
						}
						index++
						return false, nil
					})
					return err
				})
				if err != nil {
					return fmt.Errorf("file: %s: %v", sourceFile, err)
				}
			}
			return printReport(stats, os.Stdout)
		},
	}
}

func getSender(r io.Reader) (string, error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return "", err
	}
	return msg.Header.Get("From"), nil
}

func printReport(stats map[string]int, f io.Writer) error {
	for s, v := range stats {
		if _, err := fmt.Fprintf(f, "%d\t%s\n", v, s); err != nil {
			return err
		}
	}
	return nil
}
