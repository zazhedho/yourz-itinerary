package serviceshared

import (
	"time"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	domainuser "yourz-itinerary/internal/domain/user"
	"yourz-itinerary/internal/dto"
)

func TripMemberToResponse(m domaintripmember.TripMember) dto.TripMemberResponse {
	mr := dto.TripMemberResponse{
		Id:        m.Id,
		TripId:    m.TripId,
		UserId:    m.UserId,
		Role:      m.Role,
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}

	if m.UpdatedAt != nil {
		mr.UpdatedAt = new(m.UpdatedAt.Format(time.RFC3339))
	}
	if m.DeletedBy != nil {
		mr.DeletedBy = m.DeletedBy
	}
	if m.DeletedAt.Valid {
		mr.DeletedAt = new(m.DeletedAt.Time.Format(time.RFC3339))
	}

	return mr
}

func TripMemberToResponseWithUser(m domaintripmember.TripMember, user domainuser.Users) dto.TripMemberResponse {
	mr := TripMemberToResponse(m)
	mr.UserName = user.Name
	mr.UserEmail = user.Email
	mr.AvatarURL = user.AvatarURL
	return mr
}

func ItineraryDayToResponse(d domainitineraryday.ItineraryDay, items []domainitineraryitem.ItineraryItem) dto.ItineraryDayResponse {
	dr := dto.ItineraryDayResponse{
		Id:        d.Id,
		TripId:    d.TripId,
		DayNumber: d.DayNumber,
		Title:     d.Title,
		CreatedBy: d.CreatedBy,
		UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}

	if d.Date != nil {
		dr.Date = new(d.Date.Format("2006-01-02"))
	}
	if d.UpdatedAt != nil {
		dr.UpdatedAt = new(d.UpdatedAt.Format(time.RFC3339))
	}
	if d.DeletedBy != nil {
		dr.DeletedBy = d.DeletedBy
	}
	if d.DeletedAt.Valid {
		dr.DeletedAt = new(d.DeletedAt.Time.Format(time.RFC3339))
	}

	itemResponses := make([]dto.ItineraryItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, ItineraryItemToResponse(item))
	}
	dr.Items = itemResponses

	return dr
}

func ItineraryItemToResponse(item domainitineraryitem.ItineraryItem) dto.ItineraryItemResponse {
	ir := dto.ItineraryItemResponse{
		Id:           item.Id,
		DayId:        item.DayId,
		Title:        item.Title,
		Description:  item.Description,
		LocationName: item.LocationName,
		Latitude:     item.Latitude,
		Longitude:    item.Longitude,
		StartTime:    item.StartTime,
		EndTime:      item.EndTime,
		CostEstimate: item.CostEstimate,
		SortOrder:    item.SortOrder,
		CreatedBy:    item.CreatedBy,
		UpdatedBy:    item.UpdatedBy,
		CreatedAt:    item.CreatedAt.Format(time.RFC3339),
	}

	if item.UpdatedAt != nil {
		ir.UpdatedAt = new(item.UpdatedAt.Format(time.RFC3339))
	}
	if item.DeletedBy != nil {
		ir.DeletedBy = item.DeletedBy
	}
	if item.DeletedAt.Valid {
		ir.DeletedAt = new(item.DeletedAt.Time.Format(time.RFC3339))
	}

	return ir
}
