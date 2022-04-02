// Package mbox provides functions for reading and writing messages stored in mbox files
package mbox

import (
	"bufio"
	"bytes"
	"io"
)

const (
	outOfMessage = 0
	inMessage    = 1
)

// OnMessage is a callback function called on every message.
type OnMessage = func(b []byte) (bool, error)

// ReadMessages reads messages from r and calls cb for each message
func ReadMessages(r io.Reader, cb OnMessage) error {
	header := []byte("From ")
	b := bufio.NewReader(r)
	var messageBuffer bytes.Buffer
	state := outOfMessage
	for {
		if state == outOfMessage {
			peeked, err := b.Peek(len(header))
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if bytes.Equal(peeked, header) {
				state = inMessage
			}
		}
		if state == inMessage {
			sb, err := b.ReadByte()
			if err == io.EOF {
				stop, err := cb(messageBuffer.Bytes())
				if stop || err != nil {
					return err
				}
				break
			}
			if err != nil {
				return err
			}
			if sb == '\n' {
				peeked, err := b.Peek(len(header))
				if err != nil {
					if err == io.EOF {
						stop, err := cb(messageBuffer.Bytes())
						if stop || err != nil {
							return err
						}
						break
					}
					return err
				}
				if bytes.Equal(peeked, header) {
					state = outOfMessage
					stop, err := cb(messageBuffer.Bytes())
					if stop || err != nil {
						return err
					}
					messageBuffer = bytes.Buffer{}
					continue
				}
			}
			messageBuffer.WriteByte(sb)
		} else {
			// ignore
			b.ReadByte()
		}
	}
	return nil
}

// WriteMessage writes to f the message content
func WriteMessage(f io.Writer, b []byte) (int, error) {
	written := 0
	if written, err := f.Write(b); err != nil {
		return written, err
	}
	if i, err := f.Write([]byte{'\n'}); err != nil {
		return written + i, err
	} else {
		return written + i, nil
	}
}
