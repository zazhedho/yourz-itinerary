package handleritineraryitem

import (
	"errors"
	"fmt"
	"net/http"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	handlercommon "yourz-itinerary/internal/handlers/http/common"
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
	dayId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.CreateItineraryItemRequest
	handlercommon.HandleJSONMutation(ctx, handlercommon.JSONMutation[dto.CreateItineraryItemRequest, dto.ItineraryItemResponse]{
		ID:          dayId,
		Request:     &req,
		LogID:       logId,
		LogPrefix:   logPrefix,
		Operation:   "CreateItem",
		StatusCode:  http.StatusCreated,
		Message:     "Item created successfully",
		ServiceCall: h.Service.CreateItem,
		HandleError: h.handleServiceError,
	})
}

func (h *ItineraryItemHandler) UpdateItem(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[ItineraryItemHandler][UpdateItem]"
	itemId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.UpdateItineraryItemRequest
	handlercommon.HandleJSONMutation(ctx, handlercommon.JSONMutation[dto.UpdateItineraryItemRequest, dto.ItineraryItemResponse]{
		ID:          itemId,
		Request:     &req,
		LogID:       logId,
		LogPrefix:   logPrefix,
		Operation:   "UpdateItem",
		StatusCode:  http.StatusOK,
		Message:     "Item updated successfully",
		ServiceCall: h.Service.UpdateItem,
		HandleError: h.handleServiceError,
	})
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
	if !handlercommon.BindJSON(ctx, logId, &req) {
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

	if handlercommon.HandleTripAccessError(ctx, logId, err) {
		return
	}

	switch {
	case errors.Is(err, serviceitineraryitem.ErrItemNotFound):
		res := response.Response(http.StatusNotFound, "Item not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, serviceshared.ErrDayNotFound):
		res := response.Response(http.StatusNotFound, "Day not found", logId, nil)
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
