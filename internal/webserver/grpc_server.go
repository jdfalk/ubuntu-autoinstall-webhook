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

// ReportStatus handles the ReportStatus RPC.
func (s *grpcServer) ReportStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	fmt.Printf("Received status update from %s: %d%% - %s\n", req.Hostname, req.Progress, req.Message)
	return &pb.StatusResponse{Acknowledged: true}, nil
}

// StartGRPCServer starts the gRPC server on the specified port.
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}
	srv := grpc.NewServer()
	pb.RegisterInstallServiceServer(srv, &grpcServer{})
	fmt.Println("gRPC server is listening on", port)
	return srv.Serve(lis)
}
