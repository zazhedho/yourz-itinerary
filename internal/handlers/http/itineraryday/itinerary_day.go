package handleritineraryday

import (
	"errors"
	"fmt"
	"net/http"
	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	handlercommon "yourz-itinerary/internal/handlers/http/common"
	interfaceitineraryday "yourz-itinerary/internal/interfaces/itineraryday"
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
	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.CreateItineraryDayRequest
	handlercommon.HandleJSONMutation(ctx, handlercommon.JSONMutation[dto.CreateItineraryDayRequest, dto.ItineraryDayResponse]{
		ID:          tripId,
		Request:     &req,
		LogID:       logId,
		LogPrefix:   logPrefix,
		Operation:   "CreateDay",
		StatusCode:  http.StatusCreated,
		Message:     "Day created successfully",
		ServiceCall: h.Service.CreateDay,
		HandleError: h.handleServiceError,
	})
}

func (h *ItineraryDayHandler) UpdateDay(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryDayHandler][UpdateDay]"
	dayId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.UpdateItineraryDayRequest
	handlercommon.HandleJSONMutation(ctx, handlercommon.JSONMutation[dto.UpdateItineraryDayRequest, dto.ItineraryDayResponse]{
		ID:          dayId,
		Request:     &req,
		LogID:       logId,
		LogPrefix:   logPrefix,
		Operation:   "UpdateDay",
		StatusCode:  http.StatusOK,
		Message:     "Day updated successfully",
		ServiceCall: h.Service.UpdateDay,
		HandleError: h.handleServiceError,
	})
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

	if handlercommon.HandleTripAccessError(ctx, logId, err) {
		return
	}

	switch {
	case errors.Is(err, serviceshared.ErrDayNotFound):
		res := response.Response(http.StatusNotFound, "Day not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, serviceshared.ErrTripNotFound):
		res := response.Response(http.StatusNotFound, "Trip not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, serviceshared.ErrInvalidDate):
		res := response.Response(http.StatusBadRequest, "Invalid date. Must use YYYY-MM-DD.", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	default:
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
	}
}
