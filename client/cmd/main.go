package main

import (
	"fmt"
	"os"

	"github.com/Xacor/go-vault/client/internal/ui"
	"github.com/Xacor/go-vault/proto"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := proto.NewVaultServiceClient(conn)
	// secret := &proto.Secret{
	// 	Text: &proto.Text{
	// 		Data: "hello there",
	// 	},
	// }

	// // Contact the server and print out its response.
	// ctx := context.Background()
	// stream, err := c.CreateStream(ctx, &proto.Connect{User: &proto.User{Id: fmt.Sprint(rand.Int()), Name: "Simple-client"}})
	// if err != nil {
	// 	log.Fatalf("could not create stream: %v", err)
	// }
	// _, err = c.BroadcastSecret(ctx, secret)
	// if err != nil {
	// 	log.Fatalf("could not create stream: %v", err)
	// }

	// done := make(chan struct{})

	// go func() {
	// 	for {
	// 		secret, err := stream.Recv()
	// 		if err == io.EOF {
	// 			done <- struct{}{}
	// 			return
	// 		}
	// 		if err != nil {
	// 			log.Fatalf("error while receiving secret: %v", err)
	// 		}
	// 		log.Printf("Received Todo: %v", secret)
	// 	}
	// }()
	// <-done
	items := []*proto.Secret{
		{
			Name: "youtube.com",
			Credentials: &proto.Credentials{
				Login:    "nagibator1337",
				Password: "PutinLox123",
			},
		},
		{
			Name:     "Debit card",
			BankCard: &proto.BankCard{},
		},
		{
			Name: "Very secret note",
			Text: &proto.Text{
				Data: "dont tread on me",
			},
		},
		{
			Name: "youtube.com",
			Credentials: &proto.Credentials{
				Login:    "nagibator1337",
				Password: "PutinLox123",
			},
		},
		{
			Name:     "Debit card",
			BankCard: &proto.BankCard{},
		},
		{
			Name: "Very secret note",
			Text: &proto.Text{
				Data: "dont tread on me",
			},
		},
	}
	// fmt.Printf("ui.NewMainModel(items).View(): %v\n", ui.NewMainModel(items).View())

	p := tea.NewProgram(ui.InitialListModel(items))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
