package handlercommon

import (
	"fmt"

	"yourz-itinerary/internal/authscope"
	domainaudit "yourz-itinerary/internal/domain/audit"
	interfaceaudit "yourz-itinerary/internal/interfaces/audit"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
)

func WriteAudit(ctx *gin.Context, auditService interfaceaudit.ServiceAuditInterface, event domainaudit.AuditEvent, scope string) {
	if auditService == nil {
		return
	}

	scopeData := authscope.FromContext(ctx.Request.Context())
	if event.ActorUserID == "" && scopeData.ActorUserID() != "" {
		event.ActorUserID = scopeData.ActorUserID()
	}
	if event.ActorRole == "" && scopeData.ActorRole() != "" {
		event.ActorRole = scopeData.ActorRole()
	}
	event.RequestID = utils.GetRequestID(ctx)
	event.IPAddress = ctx.ClientIP()
	event.UserAgent = ctx.GetHeader("User-Agent")
	event.Metadata = utils.MergeMetadata(event.Metadata, utils.GetImpersonationMetadata(ctx))

	if err := auditService.Store(ctx.Request.Context(), event); err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelWarn, fmt.Sprintf("[%s][Audit]; failed to store audit trail: %v", scope, err))
	}
}
