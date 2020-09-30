package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/rimantoro/event_driven/profiler/entities/joblog"
	joblog_model "github.com/rimantoro/event_driven/profiler/entities/joblog/model"
	dbHelper "github.com/rimantoro/event_driven/profiler/shared/interface/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	DB         dbHelper.Database
	Collection dbHelper.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "joblogs"
)

func NewMongoRepository(DB dbHelper.Database) joblog.Repository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, jl *joblog_model.JobLog) (*joblog_model.JobLog, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, jl)
	if err != nil {
		return jl, err
	}

	return jl, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, id string) (*joblog_model.JobLog, error) {
	var (
		jl  joblog_model.JobLog
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &jl, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&jl)
	if err != nil {
		return &jl, err
	}

	return &jl, nil
}

func (m *mongoRepository) FindOneJobID(ctx context.Context, jobid string) (*joblog_model.JobLog, error) {
	var (
		jl  joblog_model.JobLog
		err error
	)

	// idHex, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return &cat, err
	// }

	err = m.Collection.FindOne(ctx, bson.M{"job_id": jobid}).Decode(&jl)
	if err != nil {
		return &jl, err
	}

	return &jl, nil
}

func (m *mongoRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]joblog_model.JobLog, int64, error) {

	var (
		joblogs []joblog_model.JobLog
		skip    int64
		opts    *options.FindOptions
	)

	skip = (p * rp) - rp
	if setsort != nil {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
			options.Find().SetSort(setsort),
		)
	} else {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
		)
	}

	cursor, err := m.Collection.Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return nil, 0, err
	}
	if cursor == nil {
		return nil, 0, fmt.Errorf("nil cursor value")
	}
	err = cursor.All(ctx, &joblogs)
	if err != nil {
		return nil, 0, err
	}

	count, err := m.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return joblogs, 0, err
	}

	return joblogs, count, err
}

func (m *mongoRepository) UpdateOne(ctx context.Context, jl *joblog_model.JobLog, id string) (*joblog_model.JobLog, error) {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return jl, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": bson.M{
		"job_id":  jl.ID,
		"status":  jl.Status,
		"message": jl.Messages,
		// "error_message": jl.Error,
		"updated_at": time.Now(),
	}}

	_, err = m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return jl, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(jl)
	if err != nil {
		return jl, err
	}
	return jl, nil
}
