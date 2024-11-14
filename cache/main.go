package main

import (
	"context"
	"fmt"
	"time"

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
	SampleCache(ctx)
}

func SampleCache(ctx context.Context) {
	c.Set(ctx, cacheKey, "10", time.Hour)
	v1, _ := c.Get(ctx, cacheKey)
	fmt.Printf("cache: value = %v\n", v1)

	c.Increment(ctx, cacheKey)
	c.Increment(ctx, cacheKey)
	c.Increment(ctx, cacheKey)

	v2, _ := c.Get(ctx, cacheKey)
	fmt.Printf("cache: value = %v\n", v2)

	c.Delete(ctx, cacheKey)
	v3, _ := c.Get(ctx, cacheKey)
	fmt.Printf("cache: value = %v\n", v3)
}
