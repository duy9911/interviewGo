package main

import (
	"context"
	"fmt"
	"interview1710/api/routers"

	"github.com/go-redis/redis/v8"
)

func main() {
	// seed.CreateTable()
	// seed.SeedData()
	// elasticDB.AddToEs()
	routers.HandleRequests()
	// ExampleClient()
}

var ctx = context.Background()

func ExampleClient() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "dadad", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
