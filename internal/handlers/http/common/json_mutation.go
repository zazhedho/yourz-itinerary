package handlercommon

import (
	"context"
	"net/http"
	"reflect"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JSONMutation[Req any, Res any] struct {
	ID          string
	Request     *Req
	LogID       uuid.UUID
	LogPrefix   string
	Operation   string
	StatusCode  int
	Message     string
	ServiceCall func(context.Context, string, string, Req) (Res, error)
	HandleError func(*gin.Context, uuid.UUID, string, error, string)
}

func HandleJSONMutation[Req any, Res any](ctx *gin.Context, mutation JSONMutation[Req, Res]) {
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	if err := ctx.BindJSON(mutation.Request); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", mutation.LogID, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(*mutation.Request), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := mutation.ServiceCall(reqCtx, scope.UserID, mutation.ID, *mutation.Request)
	if err != nil {
		mutation.HandleError(ctx, mutation.LogID, mutation.LogPrefix, err, mutation.Operation)
		return
	}

	res := response.Response(mutation.StatusCode, mutation.Message, mutation.LogID, data)
	ctx.JSON(mutation.StatusCode, res)
}
