package dto

type CreateItineraryDayRequest struct {
	Date      string `json:"date" binding:"omitempty"`
	DayNumber int    `json:"day_number" binding:"required,min=1"`
	Title     string `json:"title" binding:"omitempty,max=150"`
}

type UpdateItineraryDayRequest struct {
	Date      string `json:"date" binding:"omitempty"`
	DayNumber int    `json:"day_number" binding:"omitempty,min=1"`
	Title     string `json:"title" binding:"omitempty,max=150"`
}

type ItineraryDayResponse struct {
	Id        string                  `json:"id"`
	TripId    string                  `json:"trip_id"`
	Date      *string                 `json:"date,omitempty"`
	DayNumber int                     `json:"day_number"`
	Title     *string                 `json:"title,omitempty"`
	CreatedBy string                  `json:"created_by"`
	UpdatedBy string                  `json:"updated_by"`
	CreatedAt string                  `json:"created_at,omitempty"`
	UpdatedAt *string                 `json:"updated_at,omitempty"`
	DeletedBy *string                 `json:"deleted_by,omitempty"`
	DeletedAt *string                 `json:"deleted_at,omitempty"`
	Items     []ItineraryItemResponse `json:"items,omitempty"`
}
