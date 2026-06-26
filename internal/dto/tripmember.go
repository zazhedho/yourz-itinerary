package dto

type AddTripMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=viewer editor"`
}

type UpdateTripMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=viewer editor"`
}

type TripMemberResponse struct {
	Id        string  `json:"id"`
	TripId    string  `json:"trip_id"`
	UserId    string  `json:"user_id"`
	UserName  string  `json:"user_name,omitempty"`
	UserEmail string  `json:"user_email,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Role      string  `json:"role"`
	CreatedBy string  `json:"created_by"`
	UpdatedBy string  `json:"updated_by"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedBy *string `json:"deleted_by,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}
