package interfaceuser

import (
	"context"
	domainuser "yourz-itinerary/internal/domain/user"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoUserInterface interface {
	interfacegeneric.GenericRepository[domainuser.Users]

	GetByEmail(ctx context.Context, email string) (domainuser.Users, error)
	GetByPhone(ctx context.Context, phone string) (domainuser.Users, error)
}
