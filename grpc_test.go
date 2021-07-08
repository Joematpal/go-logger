package logger

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetCorrelationIDFromMetadata(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		need    string
		want    string
		wantErr bool
	}{
		{
			name: "should fail; no incoming metadata",
			args: args{
				ctx: context.Background(),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass; empty metadata",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{}),
			},
			need:    "some id",
			want:    "some id",
			wantErr: true,
		},
		{
			name: "should pass; with correlation_id",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("correlation_id", "with corr id")),
			},
			want: "with corr id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.need != "" {
				newID = func() string {
					return tt.need
				}
			}
			got, err := GetCorrelationIDFromMetadata(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCorrelationIDFromMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCorrelationIDFromMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

type want struct {
	Level                  string `json:"level,omitempty"`
	Msg                    string `json:"msg,omitempty"`
	CorrelationId          string `json:"correlation_id,omitempty"`
	UnaryServerInterceptor string `json:"unary_server_interceptor,omitempty"`
}

func TestLoggingStreamServerInterceptor(t *testing.T) {

	type args struct {
		ctx context.Context
	}

	newID = func() string {
		return "new_test_cor_id"
	}

	tests := []struct {
		name  string
		args  args
		wants []want
	}{
		{
			name: "should pass; with no correlation id debug msg",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{}),
			},
			wants: []want{
				{Level: "debug", Msg: "no correlation id"},
				{Level: "info", Msg: "", CorrelationId: "new_test_cor_id", UnaryServerInterceptor: "test/test/test"},
			},
		},
		{
			name: "should pass; with correlation id",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{correlationID: []string{"test_with_cor_id"}}),
			},
			wants: []want{
				{Level: "info", Msg: "", CorrelationId: "test_with_cor_id", UnaryServerInterceptor: "test/test/test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer := io.Pipe()
			logr, err := NewCorrelationLogger(
				WithEnv(dev),
				WithLevel(debug),
				WithEncoding(jsonEncoder),
				withWriter(writer),
			)

			if err != nil {
				t.Fatal(err)
			}
			var eg errgroup.Group

			eg.Go(func() error {
				scanner := bufio.NewScanner(reader)
				count := 0
				for ; scanner.Scan(); count++ {
					text := scanner.Text()
					var got want
					if err := json.Unmarshal([]byte(text), &got); err != nil {
						return err
					}

					if !reflect.DeepEqual(got, tt.wants[count]) {
						return fmt.Errorf("tt.wants[%d] failed", count)
					}
				}

				if count != len(tt.wants) {
					return fmt.Errorf("wrong length: got %d; wanted %d", count, len(tt.wants))
				}

				return scanner.Err()
			})

			middleware := LoggingStreamServerInterceptor(logr)
			ss := newServerStream(tt.args.ctx)
			info := &grpc.StreamServerInfo{
				FullMethod: "test/test/test",
			}
			handler := func(in interface{}, ss grpc.ServerStream) error {
				return nil
			}

			if err := middleware(nil, ss, info, handler); err != nil {
				t.Fatal(err)
			}

			// Close the logger writer
			writer.Close()

			if err := eg.Wait(); err != nil {
				fmt.Println("err", err)
				t.Error(err)
			}

		})
	}
}

type serverStream struct {
	context context.Context
}

func newServerStream(ctx context.Context) grpc.ServerStream {
	return &serverStream{
		context: ctx,
	}
}

func (ss *serverStream) Context() context.Context {
	if ss.context == nil {
		ss.context = context.Background()
	}
	return ss.context
}

func (ss *serverStream) RecvMsg(m interface{}) error {
	return nil
}

func (ss *serverStream) SendHeader(m metadata.MD) error {
	return nil
}

func (ss *serverStream) SendMsg(m interface{}) error {
	return nil
}

func (ss *serverStream) SetHeader(m metadata.MD) error {
	return nil
}

func (ss *serverStream) SetTrailer(m metadata.MD) {
}

func TestLoggingUnaryServerInterceptor(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	newID = func() string {
		return "new_test_cor_id"
	}

	tests := []struct {
		name  string
		args  args
		wants []want
	}{
		{
			name: "should pass; with no correlation id debug msg",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{}),
			},
			wants: []want{
				{Level: "debug", Msg: "no correlation id"},
				{Level: "info", Msg: "", CorrelationId: "new_test_cor_id", UnaryServerInterceptor: "test/test/test"},
			},
		},
		{
			name: "should pass; with correlation id",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{correlationID: []string{"test_with_cor_id"}}),
			},
			wants: []want{
				{Level: "info", Msg: "", CorrelationId: "test_with_cor_id", UnaryServerInterceptor: "test/test/test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer := io.Pipe()
			logr, err := NewCorrelationLogger(
				WithEnv(dev),
				WithLevel(debug),
				WithEncoding(jsonEncoder),
				withWriter(writer),
			)

			if err != nil {
				t.Fatal(err)
			}
			var eg errgroup.Group

			eg.Go(func() error {
				scanner := bufio.NewScanner(reader)
				count := 0
				for ; scanner.Scan(); count++ {
					text := scanner.Text()
					var got want
					if err := json.Unmarshal([]byte(text), &got); err != nil {
						return err
					}

					if !reflect.DeepEqual(got, tt.wants[count]) {
						return fmt.Errorf("tt.wants[%d] failed", count)
					}
				}

				if count != len(tt.wants) {
					return fmt.Errorf("wrong length: got %d; wanted %d", count, len(tt.wants))
				}

				return scanner.Err()
			})

			middleware := LoggingUnaryServerInterceptor(logr)
			info := &grpc.UnaryServerInfo{
				FullMethod: "test/test/test",
			}
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			}

			if _, err := middleware(tt.args.ctx, nil, info, handler); err != nil {
				t.Fatal(err)
			}

			// Close the logger writer
			writer.Close()

			if err := eg.Wait(); err != nil {
				fmt.Println("err", err)
				t.Error(err)
			}

		})
	}
}
