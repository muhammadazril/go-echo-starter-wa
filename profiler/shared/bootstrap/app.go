package bootstrap

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	redis "github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify"
	"go.uber.org/zap"

	DB "github.com/rimantoro/event_driven/profiler/shared/interface/mongo"
	"github.com/rimantoro/event_driven/profiler/shared/models"
)

var (
	App *Application
)

type Application struct {
	Logger   *zap.Logger
	Config   *viper.Viper
	Mongo    DB.Client
	Redis    *redis.Client
	Minifier struct {
		M *minify.M `json:"m"`
	} `json:"minifier"`
	Kafka struct {
		Producer *kafka.Producer
		Consumer *kafka.Consumer
	}
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Logger = InitLogger("profiler")
	App.Config = InitConfig()
	App.Mongo = InitMongo()
	App.Redis = InitRedis()
	App.Minifier.M = InitMinifier()

	kfkCfg := kafka.ConfigMap{
		"bootstrap.servers":  App.Config.GetString(models.ConfigKafkaHosts),
		"group.id":           App.Config.GetString(models.ConfigAppName),
		"enable.auto.commit": false,
		"auto.offset.reset":  "earliest",
	}
	App.Kafka.Producer = InitKafkaProducer(kfkCfg)
	App.Kafka.Consumer = InitKafkaConsumer(kfkCfg)

	InitValidator()

	os.Setenv("APP_TZ", App.Config.GetString(models.ConfigAppTimezone))
}
