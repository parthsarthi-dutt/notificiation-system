package redisclient

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// Health check
	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		log.Fatalf("❌ Redis server not running: %v", err)
	}

	log.Println("✅ Redis server connected")
}
