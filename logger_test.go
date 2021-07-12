package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
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
		name    string
		args    args
		logFile string
		want    []byte
		wantErr bool
	}{
		{
			name: "should pass; with a writer and a log file",
			args: args{
				logFile: "/tmp/test-go-logger.log",
				bytes:   bytes.NewBuffer(make([]byte, 0)),
				opts: []Option{
					WithWriters(&slowWriter{latency: 1}),
					WithLevel(debug),
					WithEncoding(jsonEncoder),
					WithEnv(dev),
				},
			},
			want: []byte("someting"),
		},
		// {
		// 	name: "should pass; with a log file",
		// 	args: args{
		// 		logFile: "/tmp/test-go-logger.log",
		// 		opts:    []Option{},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.logFile != "" {
				tt.args.opts = append(tt.args.opts, WithLogFile(tt.args.logFile))
			}
			if tt.args.bytes != nil {
				tt.args.opts = append(tt.args.opts, withWriter(tt.args.bytes))
			}
			got, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got.Debug("debug")
			got.Info("info")

			time.Sleep(time.Second)

			got.Close()

			time.Sleep(time.Second * 10)
			want, err := io.ReadAll(tt.args.bytes)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("want for real", want)
			if tt.args.logFile != "nil" {
				f, err := os.Open(tt.args.logFile)
				if err != nil {
					t.Fatal(err)
				}

				b, err := io.ReadAll(f)
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("got: %s\n", b)
				log.Printf("want: %s\n", tt.want)
				if reflect.DeepEqual(b, tt.want) || len(b) != len(tt.want) {
					t.Errorf("does not match: %s: %s", b, tt.want)
				}
			}
		})
	}
}
