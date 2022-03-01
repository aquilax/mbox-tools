package iterator

import (
	"bufio"
	"bytes"
	"io"
)

const (
	outOfMessage = 0
	inMessage    = 1
)

type OnMessage = func(b []byte) (bool, error)

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
		}
	}
	return nil
}
