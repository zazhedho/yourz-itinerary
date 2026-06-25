package handlercommon

import (
	"errors"
	"net/http"

	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleTripAccessError(ctx *gin.Context, logID uuid.UUID, err error) bool {
	switch {
	case errors.Is(err, serviceshared.ErrNotMember):
		res := response.Forbidden(logID, "You are not a member of this trip")
		ctx.JSON(http.StatusForbidden, res)
		return true
	case errors.Is(err, serviceshared.ErrAccessDenied):
		res := response.Forbidden(logID, "You do not have permission to perform this action")
		ctx.JSON(http.StatusForbidden, res)
		return true
	default:
		return false
	}
}
