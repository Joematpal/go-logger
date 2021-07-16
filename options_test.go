package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
)

func Test_writeByNewLineSync(t *testing.T) {
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
					&slowWriter{latency: time.Second}, io.Discard, &slowWriter{latency: time.Second * 2},
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

			var eg errgroup.Group

			eg.Go(func() error {
				return writeByNewLineSync(debugger, reader, tt.args.writers...)
			})

			want := &bytes.Buffer{}

			for _, msg := range tt.args.input {
				if _, err := writer.Write(msg); err != nil {
					t.Fatal(err)
				}
				if _, err := want.Write(msg); err != nil {
					t.Fatal(err)
				}
			}

			time.Sleep(time.Second * 10)

			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}

			if err := eg.Wait(); err != nil {
				t.Fatal(err)
			}

			got, err := io.ReadAll(tt.args.bytes)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(got, want.Bytes()) {
				fmt.Println("got", string(got))
				t.Errorf("diff: %v", cmp.Diff(string(got), want.String()))
			}

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

	return len(in), nil
}

type errWriter struct {
}

func (ew errWriter) Write(in []byte) (int, error) {

	return 0, errors.New("failed")
}

type testDebugger struct {
	noPrint bool
}

func (td *testDebugger) Debug(args ...interface{}) {
	if td.noPrint {
		return
	}
	log.Println(args...)

}

func (td *testDebugger) Debugf(format string, args ...interface{}) {
	if td.noPrint {
		return
	}
	log.Printf(format, args...)
	log.Println("")

}

func (td *testDebugger) Error(args ...interface{}) {
	if td.noPrint {
		return
	}
	log.Println(args...)
}

func (td *testDebugger) Errorf(format string, args ...interface{}) {
	if td.noPrint {
		return
	}
	log.Printf(format, args...)
	log.Println("")
}
