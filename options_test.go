package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

func Test_writeByNewLine(t *testing.T) {
	type args struct {
		writers []io.Writer
		input   [][]byte
		bytes   *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should pass",
			args: args{
				bytes: bytes.NewBuffer([]byte{}),
				writers: []io.Writer{
					&slowWriter{latency: time.Second}, io.Discard, &slowWriter{latency: time.Second * 5},
				},
				input: [][]byte{
					[]byte("1 this is a message\n"),
					[]byte("2 this is a message\n"),
					[]byte("3 this is a message\n"),
					[]byte("4 this is a message\n"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer := io.Pipe()
			debugger := &testDebugger{}

			tt.args.writers = append(tt.args.writers, tt.args.bytes)

			go writeByNewLine(debugger, reader, tt.args.writers...)

			for _, msg := range tt.args.input {
				if _, err := writer.Write(msg); err != nil {
					t.Fatal(err)
				}
			}

			time.Sleep(time.Second * 25)
			want, err := io.ReadAll(tt.args.bytes)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("want:", want)
		})
	}
}

type slowWriter struct {
	sync.RWMutex
	latency time.Duration
}

func (sl *slowWriter) Write(in []byte) (int, error) {
	sl.Lock()
	defer sl.Unlock()

	time.Sleep(sl.latency)
	b, err := io.ReadAll(bytes.NewReader(in))
	if err != nil {
		return len(b), err
	}
	log.Printf("slowWriter: %s", b)

	return len(in), nil
}

type testDebugger struct {
}

func (td *testDebugger) Debug(args ...interface{}) {
	log.Println(args...)
}

func (td *testDebugger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
	log.Println("")
}
