package bootstrap

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	dbHelper "github.com/rimantoro/event_driven/profiler/shared/interface/mongo"
	"github.com/rimantoro/event_driven/profiler/shared/models"
)

func InitMongo() dbHelper.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sHost := App.Config.GetString(models.ConfigMongoHost) + ":" + App.Config.GetString(models.ConfigMongoPort)

	client, _ := dbHelper.NewClient(App.Config)
	err := client.Connect(ctx)
	if err != nil {
		App.Logger.Error(fmt.Sprintf("mongodb connection error (connect) to %s", sHost), zap.String("err", err.Error()))
	}
	if err = client.Ping(ctx); err != nil {
		App.Logger.Error(fmt.Sprintf("mongodb connection error (ping) to %s", sHost), zap.String("err", err.Error()))
	} else {
		App.Logger.Info(fmt.Sprintf("Connected to MongoDB %s", sHost))
	}

	return client
}
