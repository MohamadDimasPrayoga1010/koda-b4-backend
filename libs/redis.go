package libs

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/redis/go-redis/v9"
// )

// var RedisClient *redis.Client
// var Ctx = context.Background()

// func InitRedis() *redis.Client {

// 	url := os.Getenv("REDIS_URL")
// 	password := os.Getenv("REDIS_PASSWORD")

// 	client := redis.NewClient(&redis.Options{
// 		Addr:     url,
// 		Password: password,
// 		DB:       0,
// 	})

// 	_, err := client.Ping(Ctx).Result()
// 	if err != nil {
// 		panic("Failed to connect to Redis: " + err.Error())
// 	}

// 	fmt.Println("Connected to Redis at", url)
// 	RedisClient = client
// 	return client
// }

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}

	redisURL := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")

	if redisURL == "" {
		fmt.Println("REDIS_URL not set, skipping Redis initialization")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	fmt.Println("Connected to Redis at", redisURL)
	RedisClient = client
	return client
}
