package logger

import (
	"context"
	"errors"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type wrappedStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (ws *wrappedStream) Context() context.Context {
	return ws.ctx
}

func newWrappedStream(ctx context.Context, s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{ctx, s}
}

func LoggingStreamServerInterceptor(logger Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cID, err := GetCorrelationIDFromMetadata(ss.Context())
		if err != nil {
			logger.Debug("no correlation id")
			md, _ := metadata.FromIncomingContext(ss.Context())
			md = metadata.Join(md, metadata.New(map[string]string{"correlation_id": cID}))
			ss = newWrappedStream(metadata.NewIncomingContext(ss.Context(), md), ss)
		}

		logr := logger.WithCorrelationID(cID)

		logr.Infof("stream_server_interceptor=%s", info.FullMethod)

		if err := handler(srv, ss); err != nil {
			logr.Errorf("stream_server_interceptor=%v", err)
		}
		return err
	}
}

func LoggingUnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		cID, err := GetCorrelationIDFromMetadata(ctx)
		if err != nil {
			md, _ := metadata.FromIncomingContext(ctx)
			md = metadata.Join(md, metadata.New(map[string]string{"correlation_id": cID}))
			logger.Debug("no correlation id")
			ctx = metadata.NewIncomingContext(ctx, md)
		}
		logr := logger.WithCorrelationID(cID)

		logr.Infof("unary_server_interceptor=%s request=%+v", info.FullMethod, req)

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Errorf("unary_server_interceptor=%v", err)
		}
		return resp, err
	}
}

// GetCorrelationIDFromMetadata will get the correlation_id from grpc context
// makes a new one if not present, but will still return an error
func GetCorrelationIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no incoming metdata")
	}
	values := md.Get("correlation_id")
	for _, val := range values {
		if val != "" {
			return val, nil
		}
	}

	return xid.New().String(), errors.New("empty correlation_id")
}
