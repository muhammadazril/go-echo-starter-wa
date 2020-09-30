package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocraft/work"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/rimantoro/event_driven/profiler/entities/gowa"
	gowa_model "github.com/rimantoro/event_driven/profiler/entities/gowa/model"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
	"github.com/rimantoro/event_driven/profiler/shared/models"
	"github.com/rimantoro/event_driven/profiler/shared/worker"
)

var (
	err error
)

type GowaHandler struct {
	Usecase gowa.Usecase
}

func NewHttpHandler(e *echo.Echo, uc gowa.Usecase) {
	handler := &GowaHandler{
		Usecase: uc,
	}

	g := e.Group("/wa")
	g.POST("/message", handler.SendWA)
	g.POST("/messagejob", handler.SendWAWorker)
}

// Send WA message
func (h *GowaHandler) SendWA(c echo.Context) error {

	req := new(gowa_model.WAMessage)
	err := c.Bind(req)
	if err != nil {
		return models.ReturnError(c, "52C")
	}

	var v interface{}
	if vld := bootstrap.Validate(req); len(vld) > 0 {
		bootstrap.App.Logger.Warn("error validate request SendWA", zap.Any("problem : ", vld))
		return models.ResponseValidationError(c, v, vld)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	msgId, err := h.Usecase.SendMessage(ctx, req.Number, req.Message)
	if err != nil {
		bootstrap.App.Logger.Error("error SendMessage in delivery", zap.Error(err))
		return models.ReturnError(c, "53S")
	} else {
		return c.JSON(http.StatusOK, models.ResponseTemplate{
			RC:      "00",
			Message: "success",
			Payload: gowa_model.Response{
				MsgID:  msgId,
				Status: "OK",
			},
		})
	}
}

func pushJobsSendWA(wa gowa_model.WAMessage) (*work.Job, error) {

	job := &work.Job{}
	log.Println("pushJobsSendWA")

	var i map[string]interface{}
	c, err := json.Marshal(wa)
	if err != nil {
		return job, err
	}
	json.Unmarshal(c, &i)
	job, err = worker.GetEnqueuer().Enqueue("send_wa", i)
	if err != nil {
		bootstrap.App.Logger.Error("error queue for send_wa", zap.Error(err))
		return job, err
	}

	return job, nil
}

// Send WA message via worker
func (h *GowaHandler) SendWAWorker(c echo.Context) error {

	req := new(gowa_model.WAMessage)
	err := c.Bind(req)
	if err != nil {
		return models.ReturnError(c, "52C")
	}

	var v interface{}
	if vld := bootstrap.Validate(req); len(vld) > 0 {
		bootstrap.App.Logger.Warn("error validate request SendWA", zap.Any("problem : ", vld))
		return models.ResponseValidationError(c, v, vld)
	}

	bootstrap.App.Logger.Info("Send Job WA",
		zap.Any("request", req),
	)

	job, err := pushJobsSendWA(*req)
	if err != nil {
		bootstrap.App.Logger.Error("ror push jobs SendWA", zap.Error(err))
		return models.ReturnError(c, "53S")
	}

	return c.JSON(http.StatusOK, models.ResponseTemplate{
		RC:      "00",
		Message: "success",
		Payload: gowa_model.ResponseJob{
			Message: "job processing",
			Status:  "OK",
			JobID:   job.ID,
		},
	})

}
