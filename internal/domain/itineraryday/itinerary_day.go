package domainitineraryday

import (
	"time"

	"gorm.io/gorm"
)

type ItineraryDay struct {
	Id        string         `json:"id" gorm:"column:id;primaryKey"`
	TripId    string         `json:"trip_id" gorm:"column:trip_id;not null"`
	Date      *time.Time     `json:"date,omitempty" gorm:"column:date;type:date"`
	DayNumber int            `json:"day_number" gorm:"column:day_number;not null"`
	Title     *string        `json:"title,omitempty" gorm:"column:title"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by;not null"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by;not null"`
	CreatedAt time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedBy *string        `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (ItineraryDay) TableName() string { return "itinerary_days" }
