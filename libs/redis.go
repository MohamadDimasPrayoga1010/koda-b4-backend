package libs

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() *redis.Client {


	url := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     url,       
		Password: password,   
		DB:       0,          
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	fmt.Println("Connected to Redis at", url)
	RedisClient = client
	return client
}
