package handlermenu

import (
	domainaudit "yourz-itinerary/internal/domain/audit"
	handlercommon "yourz-itinerary/internal/handlers/http/common"

	"github.com/gin-gonic/gin"
)

func (h *MenuHandler) writeAudit(ctx *gin.Context, event domainaudit.AuditEvent) {
	handlercommon.WriteAudit(ctx, h.AuditService, event, "MenuHandler")
}
