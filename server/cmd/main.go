package main

import (
	"context"
	"time"

	"github.com/Xacor/go-vault/proto"
	"github.com/Xacor/go-vault/server/internal/redis"
)

func main() {
	const url = "redis://:@localhost:6379/0"

	cli, err := redis.NewRedisClient(url)
	if err != nil {
		panic(err)
	}

	err = cli.Set(context.Background(), "text:2", proto.Text{
		Data:      "secret",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		panic(err)
	}
}
