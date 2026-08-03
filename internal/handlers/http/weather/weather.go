package handlerweather

import (
	"errors"
	"fmt"
	"net/http"

	"yourz-itinerary/internal/authscope"
	interfaceweather "yourz-itinerary/internal/interfaces/weather"
	serviceshared "yourz-itinerary/internal/services/shared"
	serviceweather "yourz-itinerary/internal/services/weather"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	Service interfaceweather.Service
}

func NewWeatherHandler(service interfaceweather.Service) *WeatherHandler {
	return &WeatherHandler{Service: service}
}

func (h *WeatherHandler) GetByDay(ctx *gin.Context) {
	logID := utils.GenerateLogId(ctx)
	logPrefix := "[WeatherHandler][GetByDay]"
	dayID, err := utils.ValidateUUID(ctx, logID)
	if err != nil {
		return
	}

	scope := authscope.FromContext(ctx.Request.Context())
	data, err := h.Service.GetByDay(ctx.Request.Context(), scope.UserID, dayID)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.GetByDay; Error: %v", logPrefix, err))
		switch {
		case errors.Is(err, serviceshared.ErrNotMember), errors.Is(err, serviceshared.ErrAccessDenied):
			ctx.JSON(http.StatusForbidden, response.Forbidden(logID, "You are not a member of this trip"))
		case errors.Is(err, serviceshared.ErrDayNotFound):
			ctx.JSON(http.StatusNotFound, response.Response(http.StatusNotFound, "Day not found", logID, nil))
		case errors.Is(err, serviceweather.ErrInternal):
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(logID))
		default:
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(logID))
		}
		return
	}

	ctx.JSON(http.StatusOK, response.Response(http.StatusOK, "Get itinerary day weather successfully", logID, data))
}
