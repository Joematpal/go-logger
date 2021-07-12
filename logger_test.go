package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_argsToString(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should pass",
			args: args{
				args: []interface{}{
					"seperate", "strings",
				},
			},
			want: "seperate strings",
		},
		{
			name: "should pass",
			args: args{
				args: []interface{}{
					"seperate", "strings", "plus one",
				},
			},
			want: "seperate strings plus one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := argsToString(tt.args.args); got != tt.want {
				t.Errorf("argsToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts    []Option
		logFile string
		bytes   *bytes.Buffer
	}
	tests := []struct {
		name       string
		args       args
		hasWriters bool
		logFile    string
		want       string
		wantErr    bool
	}{
		{
			name:       "should pass; with a writer and a log file",
			hasWriters: true,
			args: args{
				logFile: "/tmp/test-go-logger.log",
				bytes:   bytes.NewBuffer(make([]byte, 0)),
				opts: []Option{
					WithWriters(&slowWriter{latency: time.Millisecond * 500}),
					WithLevel(debug),
					WithEncoding(jsonEncoder),
					WithEnv(dev),
				},
			},
		},
		{
			name: "should pass; with a log file",
			args: args{
				logFile: "/tmp/test-go-logger.base.log",
				opts: []Option{
					WithLevel(debug),
				},
			},
			want: `[{"level":"debug","msg":"debug"},
			{"level":"info","msg":"info"}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.logFile != "" {
				tt.args.opts = append(tt.args.opts, WithLogFile(tt.args.logFile))
			}
			if tt.args.bytes != nil {
				tt.args.opts = append(tt.args.opts, withWriter(tt.args.bytes))
			}
			logr, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			logr.Debug("debug")
			logr.Info("info")

			time.Sleep(time.Second * 3)
			logr.Close()

			if tt.hasWriters {
				want, err := io.ReadAll(tt.args.bytes)
				if err != nil {
					t.Fatal(err)
				}
				f, err := os.Open(tt.args.logFile)
				if err != nil {
					t.Fatal(err)
				}

				b, err := io.ReadAll(f)
				if err != nil {
					t.Fatal(err)
				}

				if !cmp.Equal(b, want) {
					t.Errorf("does not match: %v", cmp.Diff(b, want))
				}
				return
			}

			f, err := os.Open(tt.args.logFile)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(f.Name())

			got := []testLogMsg{}
			var want []testLogMsg
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				var msg testLogMsg
				b := scanner.Bytes()

				if err := json.Unmarshal(b, &msg); err != nil {
					t.Fatal(err)
				}
				got = append(got, msg)
			}
			if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(got, want, cmpopts.IgnoreFields(testLogMsg{}, "Ts")) {
				t.Fail()
			}

		})
	}
}
