package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path"
	"strconv"

	"github.com/aquilax/mbox-tools/iterator"
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
		},
		Action: func(c *cli.Context) error {
			target := c.String("directory")
			if exists, err := afero.DirExists(fSys, target); err != nil {
				return err
			} else if !exists {
				return fmt.Errorf("target directory %s does not exist", target)
			}
			return withMBox(fSys, c.Args().Get(0), func(mb io.Reader) error {
				writers := make(map[string]io.WriteCloser)

				err := iterator.ReadMessages(mb, func(b []byte) (bool, error) {
					msg, err := mail.ReadMessage(bufio.NewReader(bytes.NewBuffer(b)))
					if err != nil {
						return true, err
					}
					time, err := msg.Header.Date()
					if err != nil {
						return false, nil
					}
					y := time.Year()
					bucketName := strconv.Itoa(y)
					if _, found := writers[bucketName]; !found {
						targetFile := path.Join(target, bucketName+".mbox")
						if exists, err := afero.Exists(fSys, targetFile); err != nil {
							return true, err
						} else if exists {
							return true, fmt.Errorf("file %s already exists", targetFile)
						}
						f, err := os.Create(targetFile)
						if err != nil {
							return true, err
						}
						writers[bucketName] = f
					}
					f := writers[bucketName]
					_, err = f.Write(b)
					if err != nil {
						return false, err
					}
					_, err = f.Write([]byte{'\n'})
					if err != nil {
						return false, err
					}
					return false, nil
				})

				for _, f := range writers {
					f.Close()
				}

				return err
			})
		},
	}
}
