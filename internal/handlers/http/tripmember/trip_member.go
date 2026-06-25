package handlertripmember

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"yourz-itinerary/internal/authscope"
	"yourz-itinerary/internal/dto"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	serviceshared "yourz-itinerary/internal/services/shared"
	servicetripmember "yourz-itinerary/internal/services/tripmember"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TripMemberHandler struct {
	Service interfacetripmember.ServiceTripMemberInterface
}

func NewTripMemberHandler(s interfacetripmember.ServiceTripMemberInterface) *TripMemberHandler {
	return &TripMemberHandler{Service: s}
}

func (h *TripMemberHandler) AddMember(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripMemberHandler][AddMember]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	var req dto.AddTripMemberRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.AddMember(reqCtx, scope.UserID, tripId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "AddMember")
		return
	}

	res := response.Response(http.StatusCreated, "Member added successfully", logId, data)
	ctx.JSON(http.StatusCreated, res)
}

func (h *TripMemberHandler) UpdateMemberRole(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripMemberHandler][UpdateMemberRole]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	memberId := ctx.Param("member_id")
	if memberId == "" {
		res := response.Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), logId, nil)
		res.Error = "member_id parameter is required"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UpdateTripMemberRoleRequest
	if err := ctx.BindJSON(&req); err != nil {
		res := response.Response(http.StatusBadRequest, "Invalid request format", logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := h.Service.UpdateMemberRole(reqCtx, scope.UserID, tripId, memberId, req)
	if err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "UpdateMemberRole")
		return
	}

	res := response.Response(http.StatusOK, "Member role updated successfully", logId, data)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripMemberHandler) RemoveMember(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripMemberHandler][RemoveMember]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	memberId := ctx.Param("member_id")
	if memberId == "" {
		res := response.Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), logId, nil)
		res.Error = "member_id parameter is required"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.Service.RemoveMember(reqCtx, scope.UserID, tripId, memberId); err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "RemoveMember")
		return
	}

	res := response.Response(http.StatusOK, "Member removed successfully", logId, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripMemberHandler) LeaveTrip(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[TripMemberHandler][LeaveTrip]"
	reqCtx := ctx.Request.Context()
	scope := authscope.FromContext(reqCtx)

	tripId, err := utils.ValidateUUID(ctx, logId)
	if err != nil {
		return
	}

	if err := h.Service.LeaveTrip(reqCtx, scope.UserID, tripId); err != nil {
		h.handleServiceError(ctx, logId, logPrefix, err, "LeaveTrip")
		return
	}

	res := response.Response(http.StatusOK, "Successfully left the trip", logId, nil)
	ctx.JSON(http.StatusOK, res)
}

func (h *TripMemberHandler) handleServiceError(ctx *gin.Context, logId uuid.UUID, logPrefix string, err error, method string) {
	logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.%s; Error: %+v", logPrefix, method, err))

	switch {
	case errors.Is(err, serviceshared.ErrNotMember):
		res := response.Forbidden(logId, "You are not a member of this trip")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, serviceshared.ErrAccessDenied):
		res := response.Forbidden(logId, "You do not have permission to perform this action")
		ctx.JSON(http.StatusForbidden, res)
	case errors.Is(err, servicetripmember.ErrMemberNotFound):
		res := response.Response(http.StatusNotFound, "Member not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, servicetripmember.ErrUserNotFound):
		res := response.Response(http.StatusNotFound, "User with this email is not registered", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, servicetripmember.ErrTripNotFound):
		res := response.Response(http.StatusNotFound, "Trip not found", logId, nil)
		ctx.JSON(http.StatusNotFound, res)
	case errors.Is(err, servicetripmember.ErrOwnerRemove):
		res := response.Response(http.StatusBadRequest, "Cannot remove the trip owner", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, servicetripmember.ErrOwnerLeave):
		res := response.Response(http.StatusBadRequest, "Owner cannot leave the trip. Transfer ownership or delete the trip instead.", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, servicetripmember.ErrOwnerRoleChange):
		res := response.Response(http.StatusBadRequest, "Cannot change the owner's role", logId, nil)
		ctx.JSON(http.StatusBadRequest, res)
	case errors.Is(err, servicetripmember.ErrDuplicateMember):
		res := response.Response(http.StatusConflict, "User is already a member of this trip", logId, nil)
		ctx.JSON(http.StatusConflict, res)
	default:
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
	}
}
