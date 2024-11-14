package client

import "github.com/go-redis/redis/v8"

type Options struct {
	Addr     string
	Password string
	DB       int
}

func NewClient(opts *Options) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
}
