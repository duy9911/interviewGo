package cache

import (
	"context"
	"encoding/json"
	"interview1710/api/models"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	return client
}

func Set(key string, ctg interface{}) {
	client := getClient()
	// serialize Post object to JSON
	json, err := json.Marshal(ctg)
	if err != nil {
		panic(err)
	}
	client.Set(ctx, key, json, 0)
}

func Get(key string) interface{} {
	client := getClient()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	post := models.Category{}
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}
	return &post
}

func Delete(key string) {
	client := getClient()
	err := client.Del(ctx, key)
	if err != nil {
		panic(err)
	}
}
