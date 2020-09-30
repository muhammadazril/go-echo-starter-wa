package worker

import (
	"context"
	"fmt"

	// "strings"

	"encoding/json"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/gocraft/work"
	redigo "github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/rimantoro/event_driven/profiler/entities/gowa"
	gowa_model "github.com/rimantoro/event_driven/profiler/entities/gowa/model"
	_gowaRepo "github.com/rimantoro/event_driven/profiler/entities/gowa/repository"
	_gowaUcase "github.com/rimantoro/event_driven/profiler/entities/gowa/usecase"
	"github.com/rimantoro/event_driven/profiler/entities/joblog"
	joblog_model "github.com/rimantoro/event_driven/profiler/entities/joblog/model"
	_joblogRepo "github.com/rimantoro/event_driven/profiler/entities/joblog/repository"
	_joblogUcase "github.com/rimantoro/event_driven/profiler/entities/joblog/usecase"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
	"github.com/rimantoro/event_driven/profiler/shared/models"
	// "github.com/bxcodec/go-clean-arch/bootstrap"
	// _catRepo "github.com/bxcodec/go-clean-arch/cat/repository/mongo"
	// "github.com/bxcodec/go-clean-arch/domain"
)

var (
	redisPool   *redigo.Pool
	enqueuer    *work.Enqueuer
	redisClient *redis.Client
	NameSpace   string
	joblogRepo  joblog.Repository
	joblogUcase joblog.Usecase
	gowaRepo    gowa.Repository
	gowaUcase   gowa.Usecase

	errMsg string
)

type Context struct {
	merchantID string
}

func init() {

	host := bootstrap.App.Config.GetString(models.ConfigRedisHost)
	port := bootstrap.App.Config.GetString(models.ConfigRedisPort)
	pass := bootstrap.App.Config.GetString(models.ConfigRedisPassword)
	db := bootstrap.App.Config.GetInt(models.ConfigRedisDatabase)

	NameSpace = bootstrap.App.Config.GetString(models.ConfigAppName)

	redisPool = &redigo.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redigo.Conn, error) {
			conn, err := redigo.Dial("tcp", fmt.Sprintf("%s:%s", host, port), redigo.DialPassword(pass), redigo.DialDatabase(db))
			if err != nil {
				bootstrap.App.Logger.Error(fmt.Sprintf("Failed to connect to Redis server %s:%s", host, port), zap.Error(err))
				// log.Print(fmt.Sprintf("Failed to connect to Redis server %s:%s", host, port))
			}
			return conn, err
		},
	}

	enqueuer = work.NewEnqueuer(NameSpace, redisPool)
}

func GetEnqueuer() *work.Enqueuer {
	return enqueuer
}

func GetRedisPool() *redigo.Pool {
	return redisPool
}

// NewWorker make a new pool workers
func NewWorker() *work.WorkerPool {

	// database := bootstrap.App.Mongo.Database(bootstrap.App.Config.GetString("mongodb.database"))
	database := bootstrap.App.Mongo.Database(bootstrap.App.Config.GetString(models.ConfigMongoDatabase))

	timeoutContext := time.Duration(30) * time.Second

	gowaRepo = _gowaRepo.NewWaRepository()
	gowaUcase = _gowaUcase.NewUsecase(gowaRepo, timeoutContext)
	joblogRepo = _joblogRepo.NewMongoRepository(database)
	joblogUcase = _joblogUcase.NewUsecase(joblogRepo, timeoutContext)

	// catRepo = _catRepo.NewMongoRepository(database)
	// catUsecase = _catUcase.NewCatUsecase(catRepo, 10*time.Second)

	// payoutRepository = PayoutRepo.NewMongoRepository(database)
	// merchantRepository = MerchantRepo.NewMongoMerchantRepository(database)
	// intUserRepository = InternalUserRepo.NewMongoInternalUserRepository(database)
	// roleRepository = RoleRepo.NewMongoRoleRepository(database)
	// permissionRepository = PermissionRepo.NewMongoRepository(database)

	pool := work.NewWorkerPool(Context{}, 10, NameSpace, redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindMerchant)

	// Map the name of jobs to handler functions
	// ExampleCustomize options:
	// pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)
	// pool.JobWithOptions("insert_cat", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).InsertCat)
	pool.JobWithOptions("send_wa", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).SendWA)
	// pool.JobWithOptions("sms_otp", work.JobOptions{Priority: 10, MaxFails: 2}, (*Context).SmsOtp)
	// pool.JobWithOptions("send_email", work.JobOptions{MaxFails: 2}, (*Context).SendEmail)
	// pool.JobWithOptions("payout_queue", work.JobOptions{MaxFails: 1}, (*Context).PayoutQueue)
	// pool.JobWithOptions("create_link", work.JobOptions{MaxFails: 1}, (*Context).CreateLink)
	// pool.JobWithOptions("callback_client", work.JobOptions{MaxFails: 5}, (*Context).SendPayoutCallbackToClient)
	// pool.JobWithOptions("disburse_approve", work.JobOptions{MaxFails: 1}, (*Context).DisburseApproval)
	// pool.JobWithOptions("disburse_reject", work.JobOptions{MaxFails: 1}, (*Context).DisburseReject)
	// pool.JobWithOptions("payout_csv", work.JobOptions{MaxFails: 1}, (*Context).PayoutFromCSV)
	// pool.JobWithOptions("get_mallaca_update", work.JobOptions{MaxFails: 1}, (*Context).GetMallacaTransaction)
	// pool.JobWithOptions("payout_reject", work.JobOptions{MaxFails: 1}, (*Context).PayoutReject)

	pool.Start()

	return pool

}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	bootstrap.App.Logger.Info("Starting Job",
		zap.String("job_name", job.Name),
		zap.String("job_id", job.ID),
		zap.Any("job_args", job.Args),
	)
	return next()
}

func (c *Context) FindMerchant(job *work.Job, next work.NextMiddlewareFunc) error {
	// If there's a merchant_id param, set it in the context for future middleware and handlers to use.
	if _, ok := job.Args["merchant_id"]; ok {
		c.merchantID = job.ArgString("merchant_id")
		if err := job.ArgError(); err != nil {
			return err
		}
	}

	return next()
}

// func (c *Context) InsertCat(job *work.Job) error {

// 	var (
// 		cat domain.Cat
// 	)

// 	arg := job.Args
// 	b, _ := json.Marshal(arg)
// 	json.Unmarshal(b, &cat)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	cat.ID = primitive.NewObjectID()
// 	cat.CreatedAt = time.Now()
// 	cat.UpdatedAt = time.Now()

// 	_, err := catRepo.InsertOne(ctx, &cat)
// 	if err != nil {
// 		log.Println("error insert cat via worker")
// 	}
// 	return nil
// }

func (c *Context) SendWA(job *work.Job) error {

	var (
		jobmsg gowa_model.WAMessage
		joblog joblog_model.JobLog
	)

	joblog.ID = primitive.NewObjectID()
	joblog.CreatedAt = time.Now()
	joblog.UpdatedAt = time.Now()
	joblog.JobID = job.ID
	joblog.Status = "sucsess"

	arg := job.Args
	b, err := json.Marshal(arg)
	if err != nil {
		bootstrap.App.Logger.Error("error marshal SendWA job args", zap.Error(err))
		joblog.JobID = job.ID
		joblog.Status = "failed"
		joblog.Error = err.Error()
	}
	json.Unmarshal(b, &jobmsg)
	joblog.Messages = jobmsg.Message
	joblog.To = jobmsg.Number
	///////////////////
	// masuk usecase repo
	//////////////////

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = gowaUcase.SendMessage(ctx, jobmsg.Number, jobmsg.Message)
	if err != nil {
		bootstrap.App.Logger.Error("error SendMessage in job SendWA", zap.Error(err))
		joblog.JobID = job.ID
		joblog.Status = "failed"
		joblog.Error = err.Error()
	}

	_, err = joblogRepo.InsertOne(ctx, &joblog)
	if err != nil {
		bootstrap.App.Logger.Error("error insert joblog via worker", zap.Error(err))
	}
	return nil
}

// func contains(slice []string, item string) bool {
// 	set := make(map[string]struct{}, len(slice))
// 	for _, s := range slice {
// 		set[s] = struct{}{}
// 	}

// 	_, ok := set[item]
// 	return ok
// }
