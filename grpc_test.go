package logger

import (
	"context"
	"fmt"
	"testing"

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
					fmt.Println("need", tt.need)
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
