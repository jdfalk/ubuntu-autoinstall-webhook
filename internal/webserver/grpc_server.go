// internal/webserver/grpc_server.go
package webserver

import (
	"context"
	"fmt"
	"net"

	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"
	"google.golang.org/grpc"
)

// grpcServer implements the InstallServiceServer interface.
type grpcServer struct {
	pb.UnimplementedInstallServiceServer
}

// ReportStatus implements the ReportStatus RPC.
func (s *grpcServer) ReportStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	fmt.Printf("Received status update from %s: %d%% - %s\n", req.Hostname, req.Progress, req.Message)
	return &pb.StatusResponse{Acknowledged: true}, nil
}

// StartGRPCServer starts a gRPC server on the provided address.
func StartGRPCServer(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	server := grpc.NewServer()
	pb.RegisterInstallServiceServer(server, &grpcServer{})
	fmt.Printf("gRPC server is listening on %s\n", address)
	return server.Serve(lis)
}
