package main

import (
	"bufio"
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"net/mail"
	"os"

	"github.com/aquilax/mbox-tools/lib/mbox"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func newMergeCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:  "merge",
		Usage: "Merges multiple mbox files",
		Subcommands: []*cli.Command{
			newMergeConcatenateCommand(fSys),
			newMergeDeduplicateCommand(fSys),
		},
	}
}

func newMergeConcatenateCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:      "concat",
		Usage:     "Prints out number of messages per sender",
		ArgsUsage: "[filename]...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output file name",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Overwrite target files",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			outputFileName := c.String("output")
			overwrite := c.Bool("overwrite")

			if !overwrite {
				if exists, err := afero.Exists(fSys, outputFileName); err != nil {
					return err
				} else if exists {
					return fmt.Errorf("file %s already exists", outputFileName)
				}
			}
			f, err := os.Create(outputFileName)
			if err != nil {
				return err
			}
			for _, sourceFile := range c.Args().Slice() {
				err := withMBox(fSys, sourceFile, func(mb io.Reader) error {
					err := mbox.ReadMessages(mb, func(b []byte) (bool, error) {
						if _, err := mbox.WriteMessage(f, b); err != nil {
							return true, err
						}
						return false, nil
					})
					return err
				})
				if err != nil {
					return fmt.Errorf("file: %s: %v", sourceFile, err)
				}
			}
			return nil
		},
	}
}

func newMergeDeduplicateCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:      "deduplicate",
		Usage:     "Prints out number of messages per sender",
		ArgsUsage: "[filename]...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output file name",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Overwrite target files",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			outputFileName := c.String("output")
			overwrite := c.Bool("overwrite")

			if !overwrite {
				if exists, err := afero.Exists(fSys, outputFileName); err != nil {
					return err
				} else if exists {
					return fmt.Errorf("file %s already exists", outputFileName)
				}
			}
			f, err := os.Create(outputFileName)
			if err != nil {
				return err
			}
			digestMap := make(map[uint64]struct{})
			for _, sourceFile := range c.Args().Slice() {
				err := withMBox(fSys, sourceFile, func(mb io.Reader) error {
					index := 0
					err := mbox.ReadMessages(mb, func(b []byte) (bool, error) {
						r := bufio.NewReader(bytes.NewBuffer(b))
						digest, err := getMessageDigest(r)
						if err != nil {
							return true, fmt.Errorf("message # %d: %v", index, err)
						}
						if _, found := digestMap[digest]; !found {
							if _, err := mbox.WriteMessage(f, b); err != nil {
								return true, err
							}
							digestMap[digest] = struct{}{}
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
			return nil
		},
	}
}

func getMessageDigest(r io.Reader) (uint64, error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return 0, err
	}

	messageId := msg.Header.Get("Message-ID")
	if messageId == "" {
		return 0, fmt.Errorf("no Message-ID header found")
	}
	h := fnv.New64a()
	_, err = h.Write([]byte(messageId))
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
