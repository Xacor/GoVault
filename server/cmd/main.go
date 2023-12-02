package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Xacor/go-vault/proto"
	"github.com/Xacor/go-vault/server/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create a new connection pool
	var conn []*service.Connection

	pool := &service.Pool{
		Connection: conn,
	}

	// Register the pool with the gRPC server
	proto.RegisterVaultServiceServer(grpcServer, pool)

	// Create a TCP listener at port 8080
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	fmt.Println("Server started at port :8080")

	// Start serving requests at port 8080
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
}
