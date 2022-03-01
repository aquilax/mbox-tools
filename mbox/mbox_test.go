package mbox

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadMessages(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		messages []string
		wantErr  bool
	}{
		{
			"works with simple message",
			`
From test@example.com Thu Jan  1 00:00:01 2020
From: test@example.com (Test)
Date: Thu, 01 Jan 2020 00:00:01 +0000
Subject: Test
This is a test.

>From Test.`,
			[]string{
				`From test@example.com Thu Jan  1 00:00:01 2020
From: test@example.com (Test)
Date: Thu, 01 Jan 2020 00:00:01 +0000
Subject: Test
This is a test.

>From Test.`,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBufferString(tt.raw)
			got := make([]string, 0)
			if err := ReadMessages(b, func(b []byte) (bool, error) { got = append(got, string(b)); return true, nil }); (err != nil) != tt.wantErr {
				t.Errorf("ReadMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.messages) {
				t.Errorf("ReadMessages() message = %v, wantErr %v", got, tt.messages)
			}
		})
	}
}

func TestWriteMessage(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantF   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bytes.Buffer{}
			got, err := WriteMessage(f, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteMessage() = %v, want %v", got, tt.want)
			}
			if gotF := f.String(); gotF != tt.wantF {
				t.Errorf("WriteMessage() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}
