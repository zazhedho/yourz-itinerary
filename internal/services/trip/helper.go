package servicetrip

import (
	"errors"
	"time"
	"unicode"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domainitineraryitem "yourz-itinerary/internal/domain/itineraryitem"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	domainuser "yourz-itinerary/internal/domain/user"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

var (
	ErrInvalidTimezone  = errors.New("invalid timezone")
	ErrInvalidCurrency  = errors.New("invalid currency code; must be a 3-letter uppercase ISO 4217 code")
	ErrInvalidDateRange = errors.New("end_date must be on or after start_date")
)

type itineraryDaySyncPlan struct {
	Create []domainitineraryday.ItineraryDay
	Update []domainitineraryday.ItineraryDay
	Delete []domainitineraryday.ItineraryDay
}

func isValidCurrencyCode(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, r := range value {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func tripToList(trip domaintrip.Trip, memberCount, dayCount int) dto.TripListResponse {
	tr := dto.TripListResponse{
		Id:           trip.Id,
		OwnerId:      trip.OwnerId,
		Title:        trip.Title,
		Destination:  trip.Destination,
		Timezone:     trip.Timezone,
		CurrencyCode: trip.CurrencyCode,
		Status:       trip.Status,
		MemberCount:  memberCount,
		DayCount:     dayCount,
		CreatedBy:    trip.CreatedBy,
		UpdatedBy:    trip.UpdatedBy,
		CreatedAt:    trip.CreatedAt.Format(time.RFC3339),
	}

	if trip.StartDate != nil {
		tr.StartDate = new(trip.StartDate.Format("2006-01-02"))
	}
	if trip.EndDate != nil {
		tr.EndDate = new(trip.EndDate.Format("2006-01-02"))
	}
	if trip.UpdatedAt != nil {
		tr.UpdatedAt = new(trip.UpdatedAt.Format(time.RFC3339))
	}
	if trip.DeletedBy != nil {
		tr.DeletedBy = trip.DeletedBy
	}
	if trip.DeletedAt.Valid {
		tr.DeletedAt = new(trip.DeletedAt.Time.Format(time.RFC3339))
	}

	return tr
}

func buildItineraryDaySyncPlan(trip domaintrip.Trip, userId string, now time.Time, existing []domainitineraryday.ItineraryDay) itineraryDaySyncPlan {
	plan := itineraryDaySyncPlan{}
	if trip.StartDate == nil {
		return plan
	}

	existingByNumber := make(map[int]domainitineraryday.ItineraryDay, len(existing))
	for _, day := range existing {
		if day.DayNumber > 0 {
			existingByNumber[day.DayNumber] = day
		}
	}

	end := *trip.StartDate
	if trip.EndDate != nil {
		end = *trip.EndDate
	}
	if end.Before(*trip.StartDate) {
		return plan
	}

	lastDayNumber := 0
	for current, dayNumber := *trip.StartDate, 1; !current.After(end); current, dayNumber = current.AddDate(0, 0, 1), dayNumber+1 {
		lastDayNumber = dayNumber
		if existingDay, exists := existingByNumber[dayNumber]; exists {
			if existingDay.Date == nil || existingDay.Date.Format("2006-01-02") != current.Format("2006-01-02") {
				existingDay.Date = new(current)
				existingDay.UpdatedBy = userId
				existingDay.UpdatedAt = &now
				plan.Update = append(plan.Update, existingDay)
			}
			continue
		}

		plan.Create = append(plan.Create, domainitineraryday.ItineraryDay{
			Id:        utils.CreateUUID(),
			TripId:    trip.Id,
			Date:      new(current),
			DayNumber: dayNumber,
			CreatedBy: userId,
			UpdatedBy: userId,
			CreatedAt: now,
		})
	}

	for _, day := range existing {
		if day.DayNumber > lastDayNumber {
			plan.Delete = append(plan.Delete, day)
		}
	}

	return plan
}

func tripToDetail(trip domaintrip.Trip, members []domaintripmember.TripMember, usersByID map[string]domainuser.Users, days []domainitineraryday.ItineraryDay, itemsByDay map[string][]domainitineraryitem.ItineraryItem) dto.TripDetailResponse {
	tr := dto.TripDetailResponse{
		Id:           trip.Id,
		OwnerId:      trip.OwnerId,
		Title:        trip.Title,
		Destination:  trip.Destination,
		Timezone:     trip.Timezone,
		CurrencyCode: trip.CurrencyCode,
		Status:       trip.Status,
		CreatedBy:    trip.CreatedBy,
		UpdatedBy:    trip.UpdatedBy,
		CreatedAt:    trip.CreatedAt.Format(time.RFC3339),
	}

	if trip.StartDate != nil {
		tr.StartDate = new(trip.StartDate.Format("2006-01-02"))
	}
	if trip.EndDate != nil {
		tr.EndDate = new(trip.EndDate.Format("2006-01-02"))
	}
	if trip.UpdatedAt != nil {
		tr.UpdatedAt = new(trip.UpdatedAt.Format(time.RFC3339))
	}
	if trip.DeletedBy != nil {
		tr.DeletedBy = trip.DeletedBy
	}
	if trip.DeletedAt.Valid {
		tr.DeletedAt = new(trip.DeletedAt.Time.Format(time.RFC3339))
	}

	memberResponses := make([]dto.TripMemberResponse, 0, len(members))
	for _, m := range members {
		if user, ok := usersByID[m.UserId]; ok {
			memberResponses = append(memberResponses, serviceshared.TripMemberToResponseWithUser(m, user))
			continue
		}
		memberResponses = append(memberResponses, serviceshared.TripMemberToResponse(m))
	}
	tr.Members = memberResponses

	dayResponses := make([]dto.ItineraryDayResponse, 0, len(days))
	for _, d := range days {
		dayResponses = append(dayResponses, serviceshared.ItineraryDayToResponse(d, itemsByDay[d.Id]))
	}
	tr.Days = dayResponses

	return tr
}
