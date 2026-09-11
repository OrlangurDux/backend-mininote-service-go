package connectors

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/redis/go-redis/v9"
	middlewares "orlangur.link/services/mini.note/handlers"
)

// RedisConnect -> Connects to Redis
func RedisConnect() *redis.Client {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     middlewares.DotEnvVariable("REDIS_URL", "localhost:6379"),
		Password: "",
		DB:       0,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("⛒ Connection Failed to Redis")
		log.Println(err)
		if err = client.Close(); err != nil {
			log.Println("Failed to close redis client:", err)
		}
	} else {
		color.Green("⛁ Connected to Redis")
	}

	return client
}
