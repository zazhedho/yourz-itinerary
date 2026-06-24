package interfaceaudit

import (
	"context"
	domainaudit "yourz-itinerary/internal/domain/audit"
	"yourz-itinerary/internal/dto"
	"yourz-itinerary/pkg/filter"
)

type ServiceAuditInterface interface {
	Store(ctx context.Context, req domainaudit.AuditEvent) error
	GetAll(ctx context.Context, params filter.BaseParams) ([]dto.AuditTrailResponse, int64, error)
	GetByID(ctx context.Context, id string) (dto.AuditTrailResponse, error)
}
