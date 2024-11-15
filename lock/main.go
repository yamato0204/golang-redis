package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yamato0204/golang-redis/client"
)

var l Locker

const key = "lock-key"

func main() {
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

	l = NewLocker(cli)
	SampleLock(ctx)
}

func Hello(i int) {
	fmt.Println("【lock】Hello, Start: ", i)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("【lock】Hello, End: ", i)
}

func SampleLock(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			DoWithLock(ctx, func() {
				Hello(i)
			})
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func DoWithLock(ctx context.Context, f func()) {
	lockInfo, err := l.TryLock(ctx, key, 5*time.Second)
	if err != nil {
		if errors.Is(err, ErrNotObtained) {
			fmt.Println("lock not obtained")
			return
		}
		fmt.Println("lock error")
		return
	}
	defer func() {
		err := lockInfo.Release(ctx)
		if err != nil && !errors.Is(err, ErrLockNotHeld) {
			fmt.Println("lock release error")
			return
		}
	}()
	f()
}
