package dto

type CreateItineraryItemRequest struct {
	Title        string   `json:"title" binding:"required,min=1,max=150"`
	Description  string   `json:"description" binding:"omitempty"`
	LocationName string   `json:"location_name" binding:"omitempty,max=200"`
	Latitude     *float64 `json:"latitude" binding:"omitempty"`
	Longitude    *float64 `json:"longitude" binding:"omitempty"`
	StartTime    string   `json:"start_time" binding:"omitempty"`
	EndTime      string   `json:"end_time" binding:"omitempty"`
	CostEstimate float64  `json:"cost_estimate" binding:"omitempty,min=0"`
	SortOrder    int      `json:"sort_order" binding:"omitempty,min=0"`
}

type UpdateItineraryItemRequest struct {
	Title        string   `json:"title" binding:"omitempty,min=1,max=150"`
	Description  string   `json:"description" binding:"omitempty"`
	LocationName string   `json:"location_name" binding:"omitempty,max=200"`
	Latitude     *float64 `json:"latitude" binding:"omitempty"`
	Longitude    *float64 `json:"longitude" binding:"omitempty"`
	StartTime    string   `json:"start_time" binding:"omitempty"`
	EndTime      string   `json:"end_time" binding:"omitempty"`
	CostEstimate float64  `json:"cost_estimate" binding:"omitempty,min=0"`
	SortOrder    int      `json:"sort_order" binding:"omitempty,min=0"`
}

type ReorderItineraryItemsRequest struct {
	ItemIds []string `json:"item_ids" binding:"required,min=1"`
}

type ItineraryItemResponse struct {
	Id           string   `json:"id"`
	DayId        string   `json:"day_id"`
	Title        string   `json:"title"`
	Description  *string  `json:"description,omitempty"`
	LocationName *string  `json:"location_name,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	StartTime    *string  `json:"start_time,omitempty"`
	EndTime      *string  `json:"end_time,omitempty"`
	CostEstimate float64  `json:"cost_estimate"`
	SortOrder    int      `json:"sort_order"`
	CreatedBy    string   `json:"created_by"`
	UpdatedBy    string   `json:"updated_by"`
	CreatedAt    string   `json:"created_at,omitempty"`
	UpdatedAt    *string  `json:"updated_at,omitempty"`
	DeletedBy    *string  `json:"deleted_by,omitempty"`
	DeletedAt    *string  `json:"deleted_at,omitempty"`
}
