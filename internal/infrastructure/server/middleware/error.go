package middleware

import (
	"context"
	"errors"
	"log"

	"homework/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewErrorMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, status.Errorf(codes.NotFound, err.Error())
			}
			if errors.Is(err, domain.ErrAlreadyExists) {
				return nil, status.Errorf(codes.AlreadyExists, err.Error())
			}
			if errors.Is(err, domain.ErrInvalidArgument) {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
			log.Printf("[interceptor.Error] method: %s; error: %s", info.FullMethod, err.Error())
			return nil, status.Error(codes.Internal, "internal server error")
		}

		return resp, err
	}
}
