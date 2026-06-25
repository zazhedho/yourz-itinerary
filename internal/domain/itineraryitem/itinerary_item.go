package domainitineraryitem

import (
	"time"

	"gorm.io/gorm"
)

type ItineraryItem struct {
	Id           string         `json:"id" gorm:"column:id;primaryKey"`
	DayId        string         `json:"day_id" gorm:"column:day_id;not null"`
	Title        string         `json:"title" gorm:"column:title;not null"`
	Description  *string        `json:"description,omitempty" gorm:"column:description"`
	LocationName *string        `json:"location_name,omitempty" gorm:"column:location_name"`
	Latitude     *float64       `json:"latitude,omitempty" gorm:"column:latitude"`
	Longitude    *float64       `json:"longitude,omitempty" gorm:"column:longitude"`
	StartTime    *string        `json:"start_time,omitempty" gorm:"column:start_time;type:time"`
	EndTime      *string        `json:"end_time,omitempty" gorm:"column:end_time;type:time"`
	CostEstimate float64        `json:"cost_estimate" gorm:"column:cost_estimate;default:0"`
	SortOrder    int            `json:"sort_order" gorm:"column:sort_order;not null"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by;not null"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by;not null"`
	CreatedAt    time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedBy    *string        `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (ItineraryItem) TableName() string { return "itinerary_items" }
