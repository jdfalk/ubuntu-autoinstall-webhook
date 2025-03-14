// internal/certadmin/auth.go
package certadmin

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor provides authentication for gRPC endpoints
type AuthInterceptor struct {
	apiKeys map[string]string // map of API keys to usernames
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(apiKeys map[string]string) *AuthInterceptor {
	return &AuthInterceptor{apiKeys: apiKeys}
}

// Unary returns a unary server interceptor function to authenticate and authorize requests
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Authenticate the request
		username, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		// Add username to the context for auditing or logging
		ctx = context.WithValue(ctx, "username", username)

		// Proceed with the request
		return handler(ctx, req)
	}
}

// authenticate validates the API key from the request metadata
func (i *AuthInterceptor) authenticate(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Get authorization header
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	// Parse "Bearer <api-key>" format
	parts := strings.SplitN(authHeader[0], " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	apiKey := parts[1]
	username, valid := i.apiKeys[apiKey]
	if !valid {
		return "", status.Error(codes.Unauthenticated, "invalid API key")
	}

	return username, nil
}
