// internal/webserver/service.go
package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// WebServer defines the interface for the webserver.
type WebServer interface {
	Start() error
	Stop() error
}

// Service implements the WebServer interface.
type Service struct{}

// NewService creates a new webserver service.
func NewService() WebServer {
	return &Service{}
}

func (s *Service) Start() error {
	// Start gRPC server in a separate goroutine.
	go func() {
		if err := StartGRPCServer(":50051"); err != nil {
			fmt.Println("Error starting gRPC server:", err)
		}
	}()

	// Wait a short time for the gRPC server to be ready.
	time.Sleep(1 * time.Second)

	// Set up grpc-gateway.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// Register the gRPC server endpoint with the gateway.
	if err := pb.RegisterInstallServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		return fmt.Errorf("failed to register gRPC gateway: %w", err)
	}

	fmt.Println("HTTP gateway is listening on :8080")
	return http.ListenAndServe(":8080", mux)
}

func (s *Service) Stop() error {
	fmt.Println("Stopping webserver")
	// TODO: implement graceful shutdown.
	return nil
}
