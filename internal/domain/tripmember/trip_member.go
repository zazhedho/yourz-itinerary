package domaintripmember

import (
	"time"

	"gorm.io/gorm"
)

type TripMember struct {
	Id        string         `json:"id" gorm:"column:id;primaryKey"`
	TripId    string         `json:"trip_id" gorm:"column:trip_id;not null"`
	UserId    string         `json:"user_id" gorm:"column:user_id;not null"`
	Role      string         `json:"role" gorm:"column:role;not null"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by;not null"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by;not null"`
	CreatedAt time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedBy *string        `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (TripMember) TableName() string { return "trip_members" }
