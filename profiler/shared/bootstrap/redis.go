package bootstrap

import (
	"fmt"

	redis "github.com/go-redis/redis/v7"
	"go.uber.org/zap"

	"github.com/rimantoro/event_driven/profiler/shared/models"
)

func InitRedis() *redis.Client {

	db := App.Config.GetInt(models.ConfigRedisDatabase)
	pass := App.Config.GetString(models.ConfigRedisPassword)
	host := App.Config.GetString(models.ConfigRedisHost)
	port := App.Config.GetString(models.ConfigRedisPort)

	address := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		App.Logger.Error(fmt.Sprintf("Failed to connect to Redis server %s:%s", host, port), zap.String("err", err.Error()))
	} else {
		redisURI := fmt.Sprintf("redis://:%s@%s:%s/%d", pass, host, port, db)
		App.Logger.Info(fmt.Sprintf("Connected to %s", redisURI))
	}

	return client
}
