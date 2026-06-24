package interfaceaudit

import (
	domainaudit "yourz-itinerary/internal/domain/audit"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoAuditInterface interface {
	interfacegeneric.GenericRepository[domainaudit.AuditTrail]
}
