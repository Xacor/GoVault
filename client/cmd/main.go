package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"

	"github.com/Xacor/go-vault/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewVaultServiceClient(conn)
	secret := &proto.Secret{
		Text: &proto.Text{
			Data: "hello there",
		},
	}

	// Contact the server and print out its response.
	ctx := context.Background()
	stream, err := c.CreateStream(ctx, &proto.Connect{User: &proto.User{Id: fmt.Sprint(rand.Int()), Name: "Simple-client"}})
	if err != nil {
		log.Fatalf("could not create stream: %v", err)
	}
	_, err = c.BroadcastSecret(ctx, secret)
	if err != nil {
		log.Fatalf("could not create stream: %v", err)
	}

	done := make(chan struct{})

	go func() {
		for {
			secret, err := stream.Recv()
			if err == io.EOF {
				done <- struct{}{}
				return
			}
			if err != nil {
				log.Fatalf("error while receiving secret: %v", err)
			}
			log.Printf("Received Todo: %v", secret)
		}
	}()
	<-done
}
