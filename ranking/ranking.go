package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Ranking interface {
	Add(ctx context.Context, key, id string, score int64) error
	Rank(ctx context.Context, key, id string) (int64, error)
	GetScore(ctx context.Context, key, id string) (int64, error)
	GetRankByScore(ctx context.Context, key string, score int64) (int64, error)
}

type redisRanking struct {
	redisClient *redis.Client
}

type RangeResult struct {
	Rank  int64
	Score int64
	ID    string
}

func NewRanking(c *redis.Client) Ranking {
	return &redisRanking{
		redisClient: c,
	}
}

func (rr *redisRanking) Add(ctx context.Context, key, id string, score int64) error {
	if cmd := rr.redisClient.ZAdd(ctx, key, &redis.Z{Score: float64(score), Member: id}); cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (rr *redisRanking) Rank(ctx context.Context, key, id string) (int64, error) {
	score, err := rr.GetScore(ctx, key, id)
	if err != nil {
		return 0, err
	}
	rank, err := rr.GetRankByScore(ctx, key, score)
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (rr *redisRanking) GetScore(ctx context.Context, key, id string) (int64, error) {
	score, err := rr.redisClient.ZScore(ctx, key, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return int64(score), nil
}

func (rr *redisRanking) GetRankByScore(ctx context.Context, key string, score int64) (int64, error) {
	count, err := rr.redisClient.ZCount(ctx, key, strconv.Itoa(int(score)+1), "+inf").Result()
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}
