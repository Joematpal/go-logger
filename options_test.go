package logger

import (
	"bytes"
	"io"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
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
					&slowWriter{latency: time.Second}, io.Discard, &slowWriter{latency: time.Second * 3},
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
				return writeByNewLine(debugger, reader, tt.args.writers...)
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
			time.Sleep(time.Second * 20)
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

			if !reflect.DeepEqual(got, want.Bytes()) {
				t.Errorf("diff: %v", cmp.Diff(got, want.Bytes()))
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
	log.Printf("slowWriter %s: %s", sl.latency, b)

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
