package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// OnConnect: func(ctx context.Context, cn *redis.Conn) error {
		// 	fmt.Println("Connected to redis server")
		// }
	})

	err := client.Set(ctx, "KeyX", "ValueIsThis", 0).Err()
	if err != nil {
		fmt.Println("Error setting value in redis : ", err)
	}

	value, err := client.Get(ctx, "KeyX").Result()
	if err != nil {
		fmt.Println("Error getting value in redis : ", err)
	} else {
		fmt.Println("Read from redis : ", value)
	}
}
