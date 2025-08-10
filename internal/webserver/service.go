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

// Service defines the WebServer interface.
type Service interface {
	Start() error
	Stop() error
}

// service implements the WebServer interface.
type service struct {
	// You can add fields for graceful shutdown, logging, etc.
}

// NewService creates and returns a new WebServer service instance.
func NewService() Service {
	return &service{}
}

// Start starts the gRPC server and HTTP gateway.
func (s *service) Start() error {
	grpcAddress := ":50051"
	// Start the gRPC server in a separate goroutine.
	go func() {
		if err := StartGRPCServer(grpcAddress); err != nil {
			fmt.Println("Error starting gRPC server:", err)
		}
	}()

	// Wait a moment for the gRPC server to start.
	time.Sleep(1 * time.Second)

	// Set up the HTTP gateway.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// Register the gRPC server endpoint with the HTTP gateway.
	if err := pb.RegisterInstallServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return fmt.Errorf("failed to register gRPC gateway: %w", err)
	}

	httpAddress := ":8080"
	fmt.Printf("HTTP gateway is listening on %s\n", httpAddress)
	return http.ListenAndServe(httpAddress, mux)
}

// Stop stops the webserver gracefully (stub implementation).
func (s *service) Stop() error {
	fmt.Println("Stopping webserver")
	// TODO: Implement graceful shutdown logic.
	return nil
}
