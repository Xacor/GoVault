package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Xacor/go-vault/proto"
)

type Connection struct {
	proto.UnimplementedVaultServiceServer
	stream proto.VaultService_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Pool struct {
	proto.UnimplementedVaultServiceServer
	Connection []*Connection
}

func (p *Pool) CreateStream(pconn *proto.Connect, stream proto.VaultService_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	log.Println("user connected: ", pconn.User.Id)
	p.Connection = append(p.Connection, conn)

	return <-conn.error
}

func (s *Pool) BroadcastSecret(ctx context.Context, secret *proto.Secret) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(secret *proto.Secret, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(secret)
				fmt.Printf("Sending message to: %v id %v", conn.id, secret.Id)

				if err != nil {
					fmt.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(secret, conn)

	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Close{}, nil
}
