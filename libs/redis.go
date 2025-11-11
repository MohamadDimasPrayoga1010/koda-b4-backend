package libs

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() *redis.Client {
	godotenv.Load()

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db := os.Getenv("REDIS_DB")

	addr := fmt.Sprintf("%s:%s", host, port)

	dbIndex := 0
	fmt.Sscanf(db, "%d", &dbIndex)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIndex,
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	RedisClient = client
	return client
}
