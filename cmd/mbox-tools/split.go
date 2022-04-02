package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path"
	"strconv"

	"github.com/aquilax/mbox-tools/lib/mbox"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func newSplitCommand(fSys afero.Fs) *cli.Command {
	return &cli.Command{
		Name:      "split",
		Usage:     "Splits an mbox file into multiple files based on message attribute",
		ArgsUsage: "[filename]...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "directory",
				Aliases:  []string{"d"},
				Usage:    "Target direcory for the split files",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "split-by",
				Aliases:  []string{"s"},
				Usage:    "Possible options: year",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "errors",
				Usage: "Target file for messages that can't be parsed",
			},
			&cli.BoolFlag{
				Name:    "overwrite",
				Aliases: []string{"o"},
				Usage:   "Overwrite target files",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			target := c.String("directory")
			splitBy := c.String("splitBy")
			overwrite := c.Bool("overwrite")
			errorsFile := c.String("errors")

			if exists, err := afero.DirExists(fSys, target); err != nil {
				return err
			} else if !exists {
				return fmt.Errorf("target directory %s does not exist", target)
			}
			for _, sourceFile := range c.Args().Slice() {
				err := withMBox(fSys, sourceFile, func(mb io.Reader) error {
					writers := make(map[string]io.WriteCloser)
					index := 0
					err := mbox.ReadMessages(mb, func(b []byte) (bool, error) {
						r := bufio.NewReader(bytes.NewBuffer(b))

						bucketName, err := getBucketName(r, splitBy, errorsFile)
						if err != nil {
							return true, fmt.Errorf("message # %d: %v", index, err)
						}
						if _, found := writers[bucketName]; !found {
							targetFile := path.Join(target, bucketName+".mbox")
							if !overwrite {
								if exists, err := afero.Exists(fSys, targetFile); err != nil {
									return true, err
								} else if exists {
									return true, fmt.Errorf("file %s already exists", targetFile)
								}
							}
							f, err := os.Create(targetFile)
							if err != nil {
								return true, err
							}
							writers[bucketName] = f
						}
						f := writers[bucketName]
						if _, err := mbox.WriteMessage(f, b); err != nil {
							return true, err
						}
						index++
						return false, nil
					})

					for _, f := range writers {
						f.Close()
					}

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

func getBucketName(r io.Reader, splitBy string, errorsFile string) (string, error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		if errorsFile != "" {
			return errorsFile, nil
		}
		return "", err
	}

	time, err := msg.Header.Date()
	if err != nil {
		if errorsFile != "" {
			return errorsFile, nil
		}
		return "", err
	}
	y := time.Year()
	return strconv.Itoa(y), nil
}
