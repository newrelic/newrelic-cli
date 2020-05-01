package main

import (
	"context"
	"time"

	pb "github.com/newrelic/newrelic-cli/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	address = "0.0.0.0:50052"
)

func main() {
	log.SetLevel(log.InfoLevel)

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
		log.Fatalf("could request execution: %v", err)
	}

	log.Infof("response: %+v", r)
}
