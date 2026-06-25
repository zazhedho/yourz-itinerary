package handleritineraryday

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
	serviceitineraryday "yourz-itinerary/internal/services/itineraryday"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItineraryDayHandler struct {
	Service interfaceitineraryday.ServiceItineraryDayInterface
}

func NewItineraryDayHandler(s interfaceitineraryday.ServiceItineraryDayInterface) *ItineraryDayHandler {
	return &ItineraryDayHandler{Service: s}
}

func (h *ItineraryDayHandler) CreateDay(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryDayHandler][CreateDay]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.CreateItineraryDayRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.CreateDay(reqCtx, scope.UserID, tripId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "CreateDay")
		return
	}

	res := response.Response(http.StatusCreated, "Day created successfully", logId, data)
	ctx.JSON(http.StatusCreated, res)
}

func (h *ItineraryDayHandler) UpdateDay(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryDayHandler][UpdateDay]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	dayId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.UpdateItineraryDayRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.UpdateDay(reqCtx, scope.UserID, dayId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "UpdateDay")
		return
	}

	res := response.Response(http.StatusOK, "Day updated successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *ItineraryDayHandler) DeleteDay(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryDayHandler][DeleteDay]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	dayId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	if err := h.Service.DeleteDay(reqCtx, scope.UserID, dayId); err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "DeleteDay")
		return
	}

	res := response.Response(http.StatusOK, "Day deleted successfully", logId, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *ItineraryDayHandler) handleServiceError(ctx *gin.Context, logId uuid.UUID, logPrefix string, err error, method string) {
	logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.%s; Error: %+v", logPrefix, method, err))

	switch {
	case errors.Is(err, serviceshared.ErrNotMember):
		res := response.Forbidden(logId, "You are not a member of this trip")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, serviceshared.ErrAccessDenied):
		res := response.Forbidden(logId, "You do not have permission to perform this action")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, serviceitineraryday.ErrDayNotFound):
		res := response.Response(http.StatusNotFound, "Day not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, serviceitineraryday.ErrTripNotFound):
		res := response.Response(http.StatusNotFound, "Trip not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	default:
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
	}
}
