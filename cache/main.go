package main

import (
	"context"

	"github.com/yamato0204/golang-redis/client"
)

var c Cache

const (
	cacheKey = "cache-key"
)

func main() {
	ctx := context.Background()
	cli := client.NewClient(&client.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		cli.Del(ctx, cacheKey)
		_ = cli.Close()

	}()
	c = NewCache(cli)
}
