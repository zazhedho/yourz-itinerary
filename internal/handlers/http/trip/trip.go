package handlertrip

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	interfacetrip "yourz-itinerary/internal/interfaces/trip"
	serviceshared "yourz-itinerary/internal/services/shared"
	servicetrip "yourz-itinerary/internal/services/trip"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TripHandler struct {
	Service interfacetrip.ServiceTripInterface
}

func NewTripHandler(s interfacetrip.ServiceTripInterface) *TripHandler {
	return &TripHandler{Service: s}
}

func (h *TripHandler) CreateTrip(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripHandler][CreateTrip]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	var req dto.CreateTripRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.CreateTrip(reqCtx, scope.UserID, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "CreateTrip")
		return
	}

	res := response.Response(http.StatusCreated, "Trip created successfully", logId, data)
	ctx.JSON(http.StatusCreated, res)
}

func (h *TripHandler) ListTrips(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripHandler][ListTrips]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	data, err := h.Service.ListTrips(reqCtx, scope.UserID)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.ListTrips; Error: %+v", logPrefix, err))
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.Response(http.StatusOK, "Trips retrieved successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripHandler) GetTripDetail(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripHandler][GetTripDetail]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	data, err := h.Service.GetTripDetail(reqCtx, scope.UserID, tripId)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "GetTripDetail")
		return
	}

	res := response.Response(http.StatusOK, "Trip detail retrieved successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripHandler) UpdateTrip(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripHandler][UpdateTrip]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.UpdateTripRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.UpdateTrip(reqCtx, scope.UserID, tripId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "UpdateTrip")
		return
	}

	res := response.Response(http.StatusOK, "Trip updated successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripHandler) DeleteTrip(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripHandler][DeleteTrip]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	if err := h.Service.DeleteTrip(reqCtx, scope.UserID, tripId); err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "DeleteTrip")
		return
	}

	res := response.Response(http.StatusOK, "Trip deleted successfully", logId, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripHandler) handleServiceError(ctx *gin.Context, logId uuid.UUID, logPrefix string, err error, method string) {
	logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.%s; Error: %+v", logPrefix, method, err))

	switch {
	case errors.Is(err, serviceshared.ErrNotMember):
		res := response.Forbidden(logId, "You are not a member of this trip")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, serviceshared.ErrAccessDenied):
		res := response.Forbidden(logId, "You do not have permission to perform this action")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, servicetrip.ErrTripNotFound):
		res := response.Response(http.StatusNotFound, "Trip not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, servicetrip.ErrInvalidTimezone):
		res := response.Response(http.StatusBadRequest, "Invalid timezone", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, servicetrip.ErrInvalidCurrency):
		res := response.Response(http.StatusBadRequest, "Invalid currency code. Must be a 3-letter uppercase ISO 4217 code.", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, servicetrip.ErrInvalidDate):
		res := response.Response(http.StatusBadRequest, "Invalid date. Must use YYYY-MM-DD.", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, servicetrip.ErrInvalidDateRange):
		res := response.Response(http.StatusBadRequest, "End date must be on or after start date", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	default:
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
	}
}
