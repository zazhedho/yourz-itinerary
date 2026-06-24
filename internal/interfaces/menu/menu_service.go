package interfacemenu

import (
	"context"
	domainmenu "yourz-itinerary/internal/domain/menu"
	"yourz-itinerary/internal/dto"
	"yourz-itinerary/pkg/filter"
)

type ServiceMenuInterface interface {
	GetByID(ctx context.Context, id string) (domainmenu.MenuItem, error)
	GetAll(ctx context.Context, params filter.BaseParams) ([]domainmenu.MenuItem, int64, error)
	GetActiveMenus(ctx context.Context) ([]domainmenu.MenuItem, error)
	GetUserMenus(ctx context.Context, userId string) ([]domainmenu.MenuItem, error)
	Update(ctx context.Context, id string, req dto.MenuUpdate) (domainmenu.MenuItem, error)
}
