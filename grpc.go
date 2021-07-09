package logger

import (
	"context"
	"errors"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	ErrNoIncomingMetadata = errors.New("no incoming metadata")
	ErrNoCorrelationID    = errors.New("no correlation id")
	newID                 = func() string {
		return xid.New().String()
	}
	Key_CorrelationID = "correlation_id"
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

func LoggingStreamServerInterceptor(logger CorrelationLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cID, err := GetCorrelationIDFromMetadata(ss.Context())
		if err != nil {
			if !errors.Is(err, ErrNoCorrelationID) {
				return err
			}
			logger.Debug("no correlation id")
			md, _ := metadata.FromIncomingContext(ss.Context())
			md = metadata.Join(md, metadata.New(map[string]string{Key_CorrelationID: cID}))
			ss = newWrappedStream(metadata.NewIncomingContext(ss.Context(), md), ss)
		}

		logr := logger.WithCorrelationID(cID)

		if l, ok := logr.(FieldLogger); ok {
			l.
				WithField("stream_server_interceptor", info.FullMethod).
				Info("")
		} else {
			logr.Infof("stream_server_interceptor=%s", info.FullMethod)
		}
		if err := handler(srv, ss); err != nil {
			logr.Errorf("stream_server_interceptor=%v", err)
		}
		return nil
	}
}

func LoggingUnaryServerInterceptor(logger CorrelationLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		cID, err := GetCorrelationIDFromMetadata(ctx)

		if err != nil {
			if !errors.Is(err, ErrNoCorrelationID) {
				return nil, err
			}
			md, _ := metadata.FromIncomingContext(ctx)
			md = metadata.Join(md, metadata.New(map[string]string{Key_CorrelationID: cID}))
			logger.Debug("no correlation id")
			ctx = metadata.NewIncomingContext(ctx, md)
		}
		logr := logger.WithCorrelationID(cID)

		if l, ok := logr.(FieldLogger); ok {
			l.
				WithFields(
					KV{"request", req},
					KV{"unary_server_interceptor", info.FullMethod},
				).
				Info("")
		} else {
			logr.Infof("unary_server_interceptor=%s request=%+v", info.FullMethod, req)
		}

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
		return "", ErrNoIncomingMetadata
	}

	values := md.Get(Key_CorrelationID)

	for _, val := range values {

		if val != "" {
			return val, nil
		}
	}

	return newID(), ErrNoCorrelationID
}
