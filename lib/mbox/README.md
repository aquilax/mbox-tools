# mbox [![Go Reference](https://pkg.go.dev/badge/github.com/aquilax/mbox-tools/lib/mbox.svg)](https://pkg.go.dev/github.com/aquilax/mbox-tools/lib/mbox)

A library for reading and writing messages from mbox files.

## Documentation

```
package mbox // import "github.com/aquilax/mbox-tools/lib/mbox"

Package mbox provides functions for reading and writing messages stored in
mbox files

FUNCTIONS

func ReadMessages(r io.Reader, cb OnMessage) error
    ReadMessages reads messages from r and calls cb for each message

func WriteMessage(f io.Writer, b []byte) (int, error)
    WriteMessage writes to f the message content


TYPES

type OnMessage = func(b []byte) (bool, error)
    OnMessage is a callback function called on every message.
```
