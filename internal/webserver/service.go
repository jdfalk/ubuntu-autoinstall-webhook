// internal/webserver/service.go
package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service implements the WebServer interface.
type Service struct{}

// NewService creates a new webserver service.
func NewService() *Service {
	return &Service{}
}

// Start launches the gRPC server and the HTTP gateway.
func (s *Service) Start() error {
	// Start gRPC server in a goroutine.
	go func() {
		if err := StartGRPCServer(":50051"); err != nil {
			fmt.Println("Error starting gRPC server:", err)
		}
	}()

	// Allow time for the gRPC server to start.
	time.Sleep(1 * time.Second)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// Register the gRPC server endpoint with the HTTP gateway.
	if err := pb.RegisterInstallServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		return fmt.Errorf("failed to register gRPC gateway: %w", err)
	}

	fmt.Println("HTTP gateway is listening on :8080")
	return http.ListenAndServe(":8080", mux)
}

// Stop is a stub for graceful shutdown.
func (s *Service) Stop() error {
	fmt.Println("Stopping webserver")
	return nil
}
