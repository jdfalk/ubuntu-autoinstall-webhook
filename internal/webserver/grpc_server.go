// internal/webserver/grpc_server.go
package webserver

import (
	"context"
	"fmt"
	"net"

	pb "github.com/jdfalk/ubuntu-autoinstall-webhook/pkg/proto"

	"google.golang.org/grpc"
)

// grpcServer implements the InstallService defined in the proto.
type grpcServer struct {
	pb.UnimplementedInstallServiceServer
}

// ReportStatus receives a status update from the client.
func (s *grpcServer) ReportStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	fmt.Printf("Received status update from %s: %d%% - %s\n", req.Hostname, req.Progress, req.Message)
	return &pb.StatusResponse{Acknowledged: true}, nil
}

// StartGRPCServer starts a gRPC server on the given port.
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterInstallServiceServer(grpcServer, &grpcServer{})
	fmt.Println("gRPC server is listening on", port)
	return grpcServer.Serve(lis)
}
