package auth

import (
	"context"
	"errors"
	"strings"

	"lucas-stellet/api-grpc-mongodb/pkg/env"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func ValidateTokenFromMeta(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	token, ok := md["authorization"]
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	if err := validateToken(token); err != nil {
		return grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	return nil
}

func validateToken(token []string) error {
	t := strings.Join(token, "")

	if t != env.Token {
		return errors.New("invalid token")
	}

	return nil
}
