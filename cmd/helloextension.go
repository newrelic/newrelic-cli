package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	pb "github.com/newrelic/newrelic-cli/rpc"
)

const (
	address = "0.0.0.0:50052"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExtensionClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ExecutionRequest{}
	r, err := c.Executions(ctx, req)

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Errorf("response: %+v", r)
}
