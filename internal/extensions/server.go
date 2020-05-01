package extensions

import (
	context "context"
	"fmt"
	"net"

	pb "github.com/newrelic/newrelic-cli/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	cmd        string
	args       []string
	grpcServer *grpc.Server
}

func NewServer(cmd string, args []string) *Server {

	var opts []grpc.ServerOption

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50052))
	if err != nil {
		log.Error(err)
	}

	grpcServer := grpc.NewServer(opts...)

	server := &Server{
		cmd:        cmd,
		args:       args,
		grpcServer: grpcServer,
	}

	pb.RegisterExtensionServer(grpcServer, server)

	go grpcServer.Serve(lis)

	return server
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) Executions(ctx context.Context, req *pb.ExecutionRequest) (*pb.CommandExecution, error) {
	c := &pb.CommandExecution{
		Command: s.cmd,
		Args:    s.args,
	}

	return c, nil
}
