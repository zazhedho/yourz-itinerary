package domaintrip

import (
	"time"

	"gorm.io/gorm"
)

type Trip struct {
	Id           string         `json:"id" gorm:"column:id;primaryKey"`
	OwnerId      string         `json:"owner_id" gorm:"column:owner_id;not null"`
	Title        string         `json:"title" gorm:"column:title;not null"`
	Destination  *string        `json:"destination,omitempty" gorm:"column:destination"`
	StartDate    *time.Time     `json:"start_date,omitempty" gorm:"column:start_date;type:date"`
	EndDate      *time.Time     `json:"end_date,omitempty" gorm:"column:end_date;type:date"`
	Timezone     string         `json:"timezone" gorm:"column:timezone;not null;default:Asia/Jakarta"`
	CurrencyCode string         `json:"currency_code" gorm:"column:currency_code;not null;default:IDR"`
	Status       string         `json:"status" gorm:"column:status;not null;default:draft"`
	CreatedBy    string         `json:"created_by" gorm:"column:created_by;not null"`
	UpdatedBy    string         `json:"updated_by" gorm:"column:updated_by;not null"`
	CreatedAt    time.Time      `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedBy    *string        `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Trip) TableName() string { return "trips" }
