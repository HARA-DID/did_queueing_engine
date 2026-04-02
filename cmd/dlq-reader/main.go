package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	streamName := os.Getenv("STREAM_NAME")
	if streamName == "" {
		streamName = "did-events"
	}
	dlqStream := streamName + ":dlq"

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse redis url: %v\n", err)
		os.Exit(1)
	}

	client := redis.NewClient(opts)
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Tailing DLQ stream: %s\n\n", dlqStream)

	lastID := "0"
	for {
		select {
		case <-ctx.Done():
			fmt.Println("exiting")
			return
		default:
		}

		streams, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{dlqStream, lastID},
			Count:   10,
			Block:   2 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			if ctx.Err() != nil {
				return
			}
			fmt.Fprintf(os.Stderr, "XREAD error: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				lastID = msg.ID
				pretty, _ := json.MarshalIndent(msg.Values, "", "  ")
				fmt.Printf("── DLQ Message ─────────────────────────\n")
				fmt.Printf("  Redis ID : %s\n", msg.ID)
				fmt.Printf("  Payload  :\n%s\n\n", pretty)
			}
		}
	}
}
