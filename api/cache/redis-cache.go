package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"interview1710/api/models"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func getClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func Set(key uint, ctg models.Category) {
	client := getClient()
	var ctx = context.Background()
	var keyS string = strconv.FormatUint(uint64(key), 10)
	// serialize Post object to JSON
	json, err := json.Marshal(ctg)
	if err != nil {
		panic(err)
	}
	client.Set(ctx, keyS, json, 0)
}

func Get(key string) *models.Category {
	var ctx = context.Background()
	fmt.Println("Get")
	client := getClient()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	post := models.Category{}
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}
	return &post
}
