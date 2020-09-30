package http

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rimantoro/event_driven/profiler/entities/joblog"
	joblog_model "github.com/rimantoro/event_driven/profiler/entities/joblog/model"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
	"github.com/rimantoro/event_driven/profiler/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	err error
)

type JobLogHandler struct {
	JobLogUsecase joblog.Usecase
}

func NewJobLogHandler(e *echo.Echo, uu joblog.Usecase) {
	handler := &JobLogHandler{
		JobLogUsecase: uu,
	}
	e.GET("/joblog/:id", handler.FindOne)
	e.GET("/joblog", handler.GetAll)
}

func (h *JobLogHandler) FindOne(c echo.Context) error {

	id := c.Param("id")

	var v interface{}
	if vld := bootstrap.VarValidate(id, "required"); len(vld) > 0 {
		bootstrap.App.Logger.Warn("error validate merchant", zap.Any("problem : ", vld))
		return models.ResponseValidationError(c, v, vld)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := h.JobLogUsecase.FindOne(ctx, id)
	if err != nil {
		bootstrap.App.Logger.Error("error find one joblog by job_id", zap.Error(err))
		return models.ReturnError(c, "52B")
	}

	return c.JSON(http.StatusOK, models.ResponseTemplate{
		RC:      "00",
		Message: "success",
		Payload: result,
	})
}

func (h *JobLogHandler) GetAll(c echo.Context) error {

	type Response struct {
		Total       int64                 `json:"total"`
		PerPage     int64                 `json:"per_page"`
		CurrentPage int64                 `json:"current_page"`
		LastPage    int64                 `json:"last_page"`
		From        int64                 `json:"from"`
		To          int64                 `json:"to"`
		JobLog      []joblog_model.JobLog `json:"joblogs"`
	}

	var (
		res   []joblog_model.JobLog
		count int64
	)

	rp, err := strconv.ParseInt(c.QueryParam("rp"), 10, 64)
	if err != nil {
		rp = 25
	}

	page, err := strconv.ParseInt(c.QueryParam("p"), 10, 64)
	if err != nil {
		page = 1
	}

	filters := bson.D{
		{"message", primitive.Regex{Pattern: ".*" + c.QueryParam("message") + ".*", Options: "i"}},
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, count, err = h.JobLogUsecase.GetAllWithPage(ctx, rp, page, filters, nil)
	if err != nil {
		return models.ReturnError(c, "53S")
	}

	result := Response{
		Total:       count,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    int64(math.Ceil(float64(count) / float64(rp))),
		From:        page*rp - rp + 1,
		To:          page * rp,
		JobLog:      res,
	}

	return c.JSON(http.StatusOK, models.ResponseTemplate{
		RC:      "00",
		Message: "success",
		Payload: result,
	})
}
