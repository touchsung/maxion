package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,     
		Password: password, 
		DB:       db,      
	})
}

func publishMessages(ctx context.Context, rdb *redis.Client, channel string, count int) {
	for i := 1; i <= count; i++ {
		message := fmt.Sprintf("message-%d", i)
		err := rdb.Publish(ctx, channel, message).Err()
		if err != nil {
			fmt.Printf("Error publishing message: %v\n", err)
			return
		}
		fmt.Printf("Published %s\n", message)
		time.Sleep(2 * time.Second)
	}
}

func receiveMessages(ctx context.Context, subscriber *redis.PubSub, count int) {
	fmt.Println("Starting to receive messages...")
	for i := 1; i <= count; i++ {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			return
		}
		fmt.Printf("Received message from channel %s: %s\n", msg.Channel, msg.Payload)
	}
}

func main() {
	const (
		channel     = "mychannel"
		messageCount = 5
	)

	ctx := context.Background()
	rdb := NewRedisClient("localhost:6379", "", 0)
	defer rdb.Close()

	subscriber := rdb.Subscribe(ctx, channel)
	defer subscriber.Close()

	go publishMessages(ctx, rdb, channel, messageCount)

	receiveMessages(ctx, subscriber, messageCount)
}
