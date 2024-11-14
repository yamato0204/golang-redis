package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yamato0204/golang-redis/client"
)

var r Ranking

const key = "key-ranking-simple"

func main() {
	fmt.Println("start")
	ctx := context.Background()
	cli := client.NewClient(&client.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer func() {
		cli.Del(ctx, key)
		_ = cli.Close()
	}()

	r = NewRanking(cli)

	Sample(ctx)

}

func Sample(ctx context.Context) {
	for _, user := range []struct {
		ID        string
		Score     int64
		TimeStamp time.Time
	}{
		{ID: "user1", Score: 100},
		{ID: "user2", Score: 200},
		{ID: "user3", Score: 300},
		{ID: "user4", Score: 300},
		{ID: "user5", Score: 500},
		{ID: "user6", Score: 600},
		{ID: "user7", Score: 600},
		{ID: "user8", Score: 600},
		{ID: "user9", Score: 900},
		{ID: "user10", Score: 1000},
	} {
		r.Add(ctx, key, user.ID, user.Score)
	}

	rankResult, _ := r.Rank(ctx, key, "user6")
	fmt.Printf("【ranking】rankResult(user6) = %#v\n", rankResult)
	rankResult, _ = r.Rank(ctx, key, "user7")
	fmt.Printf("【ranking】rankResult(user7) = %#v\n", rankResult)
	rankResult, _ = r.Rank(ctx, key, "user8")
	fmt.Printf("【ranking】rankResult(user8) = %#v\n", rankResult)

	rankByScoreResult, _ := r.GetRankByScore(ctx, key, 600)
	fmt.Printf("【ranking】rankByScoreResult = %#v\n", rankByScoreResult)

	scoreResult, _ := r.GetScore(ctx, key, "user7")
	fmt.Printf("【ranking】scoreResult = %#v\n", scoreResult)

}
