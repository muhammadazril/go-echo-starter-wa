package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/rimantoro/event_driven/profiler/entities/client"
	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
	"github.com/rimantoro/event_driven/profiler/shared/models"
)

var (
	err error
)

type ClientProducerHandler struct {
	Usecase client.Usecase
}

func NewHttpHandler(e *echo.Echo, uc client.Usecase) {
	handler := &ClientProducerHandler{
		Usecase: uc,
	}

	g := e.Group("/test")
	g.POST("/client", handler.Store)
}

func (h *ClientProducerHandler) Store(c echo.Context) error {

	type reqStruct struct {
		ClientID string `json:"client_id" validate:"required,min=8,alphanum"`
	}

	req := new(reqStruct)
	if err := c.Bind(req); err != nil {
		return models.ReturnError(c, "52C")
	}

	var v interface{}
	if vld := bootstrap.Validate(req); len(vld) > 0 {
		bootstrap.App.Logger.Warn("error validate merchant", zap.Any("problem : ", vld))
		return models.ResponseValidationError(c, v, vld)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// produce message
	if err = h.Usecase.PostStreamRegisterClient(ctx, req.ClientID); err != nil {
		return models.ReturnError(c, "52C")
	}

	return c.JSON(http.StatusOK, models.ResponseTemplate{
		RC:      "00",
		Message: "success",
		Payload: "OK",
	})
}
