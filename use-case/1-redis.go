package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,     
		Password: password, 
		DB:       db,      
	})
}

func SetKey(client *redis.Client, key, value string) error {
	return client.Set(ctx, key, value, 0).Err()
}

func GetKey(client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func main() {
	client := NewRedisClient("localhost:6379", "", 0)
	defer client.Close()

	if err := SetKey(client, "greeting", "Hello, Redis!"); err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}

	val, err := GetKey(client, "greeting")
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}

	fmt.Println("Value from Redis:", val)
}
