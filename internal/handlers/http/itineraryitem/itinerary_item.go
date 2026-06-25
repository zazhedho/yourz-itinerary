package handleritineraryitem

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	interfaceitineraryitem "yourz-itinerary/internal/interfaces/itineraryitem"
	serviceitineraryitem "yourz-itinerary/internal/services/itineraryitem"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItineraryItemHandler struct {
	Service interfaceitineraryitem.ServiceItineraryItemInterface
}

func NewItineraryItemHandler(s interfaceitineraryitem.ServiceItineraryItemInterface) *ItineraryItemHandler {
	return &ItineraryItemHandler{Service: s}
}

func (h *ItineraryItemHandler) CreateItem(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryItemHandler][CreateItem]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	dayId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.CreateItineraryItemRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.CreateItem(reqCtx, scope.UserID, dayId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "CreateItem")
		return
	}

	res := response.Response(http.StatusCreated, "Item created successfully", logId, data)
	ctx.JSON(http.StatusCreated, res)
}

func (h *ItineraryItemHandler) UpdateItem(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryItemHandler][UpdateItem]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	itemId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.UpdateItineraryItemRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.UpdateItem(reqCtx, scope.UserID, itemId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "UpdateItem")
		return
	}

	res := response.Response(http.StatusOK, "Item updated successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *ItineraryItemHandler) DeleteItem(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryItemHandler][DeleteItem]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	itemId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	if err := h.Service.DeleteItem(reqCtx, scope.UserID, itemId); err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "DeleteItem")
		return
	}

	res := response.Response(http.StatusOK, "Item deleted successfully", logId, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *ItineraryItemHandler) ReorderItems(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryItemHandler][ReorderItems]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	dayId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.ReorderItineraryItemsRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.Service.ReorderItems(reqCtx, scope.UserID, dayId, req); err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "ReorderItems")
		return
	}

	res := response.Response(http.StatusOK, "Items reordered successfully", logId, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *ItineraryItemHandler) handleServiceError(ctx *gin.Context, logId uuid.UUID, logPrefix string, err error, method string) {
	logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.%s; Error: %+v", logPrefix, method, err))

	switch {
	case errors.Is(err, serviceshared.ErrNotMember):
		res := response.Forbidden(logId, "You are not a member of this trip")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, serviceshared.ErrAccessDenied):
		res := response.Forbidden(logId, "You do not have permission to perform this action")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, serviceitineraryitem.ErrItemNotFound):
		res := response.Response(http.StatusNotFound, "Item not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, serviceitineraryitem.ErrInvalidTime):
		res := response.Response(http.StatusBadRequest, "Invalid time. Must use HH:MM.", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, serviceitineraryitem.ErrInvalidCoordinates):
		res := response.Response(http.StatusBadRequest, "Latitude and longitude must both be provided", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, serviceitineraryitem.ErrInvalidLatitude):
		res := response.Response(http.StatusBadRequest, "Latitude must be between -90 and 90", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, serviceitineraryitem.ErrInvalidLongitude):
		res := response.Response(http.StatusBadRequest, "Longitude must be between -180 and 180", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, serviceitineraryitem.ErrReorderDifferentDay):
		res := response.Response(http.StatusBadRequest, "All items must belong to the same day", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, serviceitineraryitem.ErrReorderEmpty):
		res := response.Response(http.StatusBadRequest, "item_ids must not be empty", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, serviceitineraryitem.ErrReorderItemsNotFound):
		res := response.Response(http.StatusBadRequest, "One or more items not found", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	default:
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
	}
}
