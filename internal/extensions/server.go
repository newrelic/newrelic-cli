package extensions

import (
	context "context"
	"fmt"
	"net"

	pb "github.com/newrelic/newrelic-cli/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct{}

func NewServer() *Server {
	var server *Server
	var opts []grpc.ServerOption

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50052))
	if err != nil {
		log.Error(err)
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterExtensionServer(grpcServer, server)

	grpcServer.Serve(lis)

	return server
}

func (s *Server) Executions(ctx context.Context, req *pb.ExecutionRequest) (*pb.CommandExecution, error) {
	c := &pb.CommandExecution{
		Command: "hello",
		Args: []string{
			"world",
		},
	}

	return c, nil
}
