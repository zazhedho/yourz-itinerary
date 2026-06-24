package interfaceappconfig

import (
	"context"
	domainappconfig "yourz-itinerary/internal/domain/appconfig"
	interfacegeneric "yourz-itinerary/internal/interfaces/generic"
)

type RepoAppConfigInterface interface {
	interfacegeneric.GenericRepository[domainappconfig.AppConfig]

	GetByKey(ctx context.Context, configKey string) (domainappconfig.AppConfig, error)
}
