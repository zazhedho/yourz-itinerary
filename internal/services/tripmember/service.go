package servicetripmember

import (
	"context"
	"strings"
	"time"

	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	"yourz-itinerary/internal/dto"
	interfacetrip "yourz-itinerary/internal/interfaces/trip"
	interfacetripmember "yourz-itinerary/internal/interfaces/tripmember"
	interfaceuser "yourz-itinerary/internal/interfaces/user"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/utils"
)

type TripMemberService struct {
	tripRepo   interfacetrip.RepoTripInterface
	memberRepo interfacetripmember.RepoTripMemberInterface
	userRepo   interfaceuser.RepoUserInterface
}

func NewTripMemberService(
	tripRepo interfacetrip.RepoTripInterface,
	memberRepo interfacetripmember.RepoTripMemberInterface,
	userRepo interfaceuser.RepoUserInterface,
) *TripMemberService {
	return &TripMemberService{
		tripRepo:   tripRepo,
		memberRepo: memberRepo,
		userRepo:   userRepo,
	}
}

func (s *TripMemberService) AddMember(ctx context.Context, userId, tripId string, req dto.AddTripMemberRequest) (dto.TripMemberResponse, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return dto.TripMemberResponse{}, ErrTripNotFound
	}

	if trip.OwnerId != userId {
		return dto.TripMemberResponse{}, serviceshared.ErrAccessDenied
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	targetUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || targetUser.Id == "" {
		return dto.TripMemberResponse{}, ErrUserNotFound
	}

	existing, _ := s.memberRepo.GetActiveByTripAndUser(ctx, tripId, targetUser.Id)
	if existing.Id != "" {
		return dto.TripMemberResponse{}, ErrDuplicateMember
	}

	role := strings.TrimSpace(req.Role)
	if role != serviceshared.TripRoleViewer && role != serviceshared.TripRoleEditor {
		role = serviceshared.TripRoleViewer
	}

	now := time.Now()
	member := domaintripmember.TripMember{
		Id:        utils.CreateUUID(),
		TripId:    tripId,
		UserId:    targetUser.Id,
		Role:      role,
		CreatedBy: userId,
		UpdatedBy: userId,
		CreatedAt: now,
	}

	if err := s.memberRepo.Store(ctx, member); err != nil {
		return dto.TripMemberResponse{}, err
	}

	return memberToResponse(member), nil
}

func (s *TripMemberService) UpdateMemberRole(ctx context.Context, userId, tripId, memberId string, req dto.UpdateTripMemberRoleRequest) (dto.TripMemberResponse, error) {
	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return dto.TripMemberResponse{}, ErrTripNotFound
	}

	if trip.OwnerId != userId {
		return dto.TripMemberResponse{}, serviceshared.ErrAccessDenied
	}

	member, err := s.memberRepo.GetByID(ctx, memberId)
	if err != nil {
		return dto.TripMemberResponse{}, ErrMemberNotFound
	}

	if member.TripId != tripId {
		return dto.TripMemberResponse{}, ErrMemberNotFound
	}

	if member.Role == serviceshared.TripRoleOwner {
		return dto.TripMemberResponse{}, ErrOwnerRoleChange
	}

	role := strings.TrimSpace(req.Role)
	if role != serviceshared.TripRoleViewer && role != serviceshared.TripRoleEditor {
		return dto.TripMemberResponse{}, serviceshared.ErrAccessDenied
	}

	member.Role = role
	member.UpdatedBy = userId
	member.UpdatedAt = new(time.Now())

	if err := s.memberRepo.Update(ctx, member); err != nil {
		return dto.TripMemberResponse{}, err
	}

	return memberToResponse(member), nil
}

func (s *TripMemberService) RemoveMember(ctx context.Context, userId, tripId, memberId string) error {
	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return ErrTripNotFound
	}

	if trip.OwnerId != userId {
		return serviceshared.ErrAccessDenied
	}

	member, err := s.memberRepo.GetByID(ctx, memberId)
	if err != nil {
		return ErrMemberNotFound
	}

	if member.TripId != tripId {
		return ErrMemberNotFound
	}

	if member.UserId == trip.OwnerId {
		return ErrOwnerRemove
	}

	return s.memberRepo.SoftDelete(ctx, member.Id, userId)
}

func (s *TripMemberService) LeaveTrip(ctx context.Context, userId, tripId string) error {
	trip, err := s.tripRepo.GetByID(ctx, tripId)
	if err != nil {
		return ErrTripNotFound
	}

	if trip.OwnerId == userId {
		return ErrOwnerLeave
	}

	member, err := s.memberRepo.GetActiveByTripAndUser(ctx, tripId, userId)
	if err != nil {
		return serviceshared.ErrNotMember
	}

	return s.memberRepo.SoftDelete(ctx, member.Id, userId)
}
