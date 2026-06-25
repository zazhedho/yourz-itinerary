package dto

type CreateTripRequest struct {
	Title        string `json:"title" binding:"required,min=1,max=150"`
	Destination  string `json:"destination" binding:"omitempty,max=150"`
	StartDate    string `json:"start_date" binding:"omitempty"`
	EndDate      string `json:"end_date" binding:"omitempty"`
	Timezone     string `json:"timezone" binding:"omitempty,max=64"`
	CurrencyCode string `json:"currency_code" binding:"omitempty,max=3"`
}

type UpdateTripRequest struct {
	Title        string `json:"title" binding:"omitempty,min=1,max=150"`
	Destination  string `json:"destination" binding:"omitempty,max=150"`
	StartDate    string `json:"start_date" binding:"omitempty"`
	EndDate      string `json:"end_date" binding:"omitempty"`
	Timezone     string `json:"timezone" binding:"omitempty,max=64"`
	CurrencyCode string `json:"currency_code" binding:"omitempty,max=3"`
	Status       string `json:"status" binding:"omitempty,max=30"`
}

type TripListResponse struct {
	Id           string  `json:"id"`
	OwnerId      string  `json:"owner_id"`
	Title        string  `json:"title"`
	Destination  *string `json:"destination,omitempty"`
	StartDate    *string `json:"start_date,omitempty"`
	EndDate      *string `json:"end_date,omitempty"`
	Timezone     string  `json:"timezone"`
	CurrencyCode string  `json:"currency_code"`
	Status       string  `json:"status"`
	MemberCount  int     `json:"member_count"`
	DayCount     int     `json:"day_count"`
	CreatedBy    string  `json:"created_by"`
	UpdatedBy    string  `json:"updated_by"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
	DeletedBy    *string `json:"deleted_by,omitempty"`
	DeletedAt    *string `json:"deleted_at,omitempty"`
}

type TripDetailResponse struct {
	Id           string                 `json:"id"`
	OwnerId      string                 `json:"owner_id"`
	Title        string                 `json:"title"`
	Destination  *string                `json:"destination,omitempty"`
	StartDate    *string                `json:"start_date,omitempty"`
	EndDate      *string                `json:"end_date,omitempty"`
	Timezone     string                 `json:"timezone"`
	CurrencyCode string                 `json:"currency_code"`
	Status       string                 `json:"status"`
	CreatedBy    string                 `json:"created_by"`
	UpdatedBy    string                 `json:"updated_by"`
	CreatedAt    string                 `json:"created_at,omitempty"`
	UpdatedAt    *string                `json:"updated_at,omitempty"`
	DeletedBy    *string                `json:"deleted_by,omitempty"`
	DeletedAt    *string                `json:"deleted_at,omitempty"`
	Members      []TripMemberResponse   `json:"members"`
	Days         []ItineraryDayResponse `json:"days"`
}
