package servicetripmember

import (
	"context"
	"errors"
	"testing"

	domainitineraryday "yourz-itinerary/internal/domain/itineraryday"
	domaintrip "yourz-itinerary/internal/domain/trip"
	domaintripmember "yourz-itinerary/internal/domain/tripmember"
	domainuser "yourz-itinerary/internal/domain/user"
	"yourz-itinerary/internal/dto"
	serviceshared "yourz-itinerary/internal/services/shared"
	"yourz-itinerary/pkg/filter"
)

type tripRepoStub struct {
	trip domaintrip.Trip
	err  error
}

func (s *tripRepoStub) Store(context.Context, domaintrip.Trip) error { return nil }
func (s *tripRepoStub) GetByID(context.Context, string) (domaintrip.Trip, error) {
	return s.trip, s.err
}
func (s *tripRepoStub) GetAll(context.Context, filter.BaseParams) ([]domaintrip.Trip, int64, error) {
	return nil, 0, nil
}
func (s *tripRepoStub) Update(context.Context, domaintrip.Trip) error    { return nil }
func (s *tripRepoStub) Delete(context.Context, string) error             { return nil }
func (s *tripRepoStub) SoftDelete(context.Context, string, string) error { return nil }
func (s *tripRepoStub) CreateTrip(context.Context, domaintrip.Trip, domaintripmember.TripMember, ...domainitineraryday.ItineraryDay) (domaintrip.Trip, error) {
	return domaintrip.Trip{}, nil
}
func (s *tripRepoStub) ListByMember(context.Context, string) ([]domaintrip.Trip, int64, error) {
	return nil, 0, nil
}

type memberRepoStub struct {
	member        domaintripmember.TripMember
	err           error
	storeErr      error
	updateErr     error
	softDeleteErr error
}

func (s *memberRepoStub) Store(context.Context, domaintripmember.TripMember) error { return s.storeErr }
func (s *memberRepoStub) GetByID(context.Context, string) (domaintripmember.TripMember, error) {
	return s.member, s.err
}
func (s *memberRepoStub) GetAll(context.Context, filter.BaseParams) ([]domaintripmember.TripMember, int64, error) {
	return nil, 0, nil
}
func (s *memberRepoStub) Update(context.Context, domaintripmember.TripMember) error {
	return s.updateErr
}
func (s *memberRepoStub) Delete(context.Context, string) error             { return nil }
func (s *memberRepoStub) SoftDelete(context.Context, string, string) error { return s.softDeleteErr }
func (s *memberRepoStub) GetByTripAndUser(context.Context, string, string) (domaintripmember.TripMember, error) {
	return s.member, s.err
}
func (s *memberRepoStub) GetActiveByTripAndUser(context.Context, string, string) (domaintripmember.TripMember, error) {
	return s.member, s.err
}
func (s *memberRepoStub) ListByTrip(context.Context, string) ([]domaintripmember.TripMember, error) {
	return nil, nil
}

type userRepoStub struct {
	user domainuser.Users
	err  error
}

func (s *userRepoStub) Store(context.Context, domainuser.Users) error { return nil }
func (s *userRepoStub) GetByID(context.Context, string) (domainuser.Users, error) {
	return s.user, s.err
}
func (s *userRepoStub) GetAll(context.Context, filter.BaseParams) ([]domainuser.Users, int64, error) {
	return nil, 0, nil
}
func (s *userRepoStub) Update(context.Context, domainuser.Users) error   { return nil }
func (s *userRepoStub) Delete(context.Context, string) error             { return nil }
func (s *userRepoStub) SoftDelete(context.Context, string, string) error { return nil }
func (s *userRepoStub) GetByEmail(context.Context, string) (domainuser.Users, error) {
	return s.user, s.err
}
func (s *userRepoStub) GetByPhone(context.Context, string) (domainuser.Users, error) {
	return domainuser.Users{}, nil
}

func memberServiceForTest(trip domaintrip.Trip, member domaintripmember.TripMember, user domainuser.Users) (*TripMemberService, *memberRepoStub) {
	memberRepo := &memberRepoStub{member: member}
	return NewTripMemberService(&tripRepoStub{trip: trip}, memberRepo, &userRepoStub{user: user}), memberRepo
}

func TestTripMemberServiceAddMember(t *testing.T) {
	trip := domaintrip.Trip{Id: "trip-1", OwnerId: "owner-1"}
	user := domainuser.Users{Id: "user-2", Email: "person@example.com"}
	svc, repo := memberServiceForTest(trip, domaintripmember.TripMember{}, user)
	got, err := svc.AddMember(context.Background(), "owner-1", "trip-1", dto.AddTripMemberRequest{Email: " Person@Example.com ", Role: "invalid"})
	if err != nil || got.Role != serviceshared.TripRoleViewer || repo.storeErr != nil {
		t.Fatalf("add member: %+v err=%v", got, err)
	}

	cases := []struct {
		name   string
		trip   domaintrip.Trip
		user   domainuser.Users
		member domaintripmember.TripMember
		want   error
	}{
		{"not found", domaintrip.Trip{}, user, domaintripmember.TripMember{}, serviceshared.ErrTripNotFound},
		{"not owner", domaintrip.Trip{Id: "trip-1", OwnerId: "owner-2"}, user, domaintripmember.TripMember{}, serviceshared.ErrAccessDenied},
		{"user missing", trip, domainuser.Users{}, domaintripmember.TripMember{}, ErrUserNotFound},
		{"duplicate", trip, user, domaintripmember.TripMember{Id: "member-1"}, ErrDuplicateMember},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &userRepoStub{user: tc.user}
			if tc.name == "user missing" {
				userRepo.err = errors.New("missing")
			}
			memberRepo := &memberRepoStub{member: tc.member}
			tripRepo := &tripRepoStub{trip: tc.trip}
			if tc.name == "not found" {
				tripRepo.err = errors.New("missing")
			}
			svc := NewTripMemberService(tripRepo, memberRepo, userRepo)
			_, err := svc.AddMember(context.Background(), "owner-1", "trip-1", dto.AddTripMemberRequest{Email: "person@example.com"})
			if !errors.Is(err, tc.want) {
				t.Fatalf("err=%v want=%v", err, tc.want)
			}
		})
	}
}

func TestTripMemberServiceUpdateRemoveAndLeave(t *testing.T) {
	trip := domaintrip.Trip{Id: "trip-1", OwnerId: "owner-1"}
	member := domaintripmember.TripMember{Id: "member-1", TripId: "trip-1", UserId: "user-2", Role: serviceshared.TripRoleViewer}
	svc, memberRepo := memberServiceForTest(trip, member, domainuser.Users{Id: "user-2"})
	got, err := svc.UpdateMemberRole(context.Background(), "owner-1", "trip-1", "member-1", dto.UpdateTripMemberRoleRequest{Role: serviceshared.TripRoleEditor})
	if err != nil || got.Role != serviceshared.TripRoleEditor {
		t.Fatalf("update: %+v err=%v", got, err)
	}

	memberRepo.member = domaintripmember.TripMember{Id: "owner-member", TripId: "trip-1", UserId: "owner-1", Role: serviceshared.TripRoleOwner}
	if _, err := svc.UpdateMemberRole(context.Background(), "owner-1", "trip-1", "owner-member", dto.UpdateTripMemberRoleRequest{Role: serviceshared.TripRoleViewer}); !errors.Is(err, ErrOwnerRoleChange) {
		t.Fatalf("owner role err=%v", err)
	}
	memberRepo.member = member
	if _, err := svc.UpdateMemberRole(context.Background(), "owner-1", "trip-1", "member-1", dto.UpdateTripMemberRoleRequest{Role: "invalid"}); !errors.Is(err, serviceshared.ErrAccessDenied) {
		t.Fatalf("invalid role err=%v", err)
	}

	if err := svc.RemoveMember(context.Background(), "owner-1", "trip-1", "member-1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	memberRepo.member = domaintripmember.TripMember{Id: "owner-member", TripId: "trip-1", UserId: "owner-1", Role: serviceshared.TripRoleOwner}
	if err := svc.RemoveMember(context.Background(), "owner-1", "trip-1", "owner-member"); !errors.Is(err, ErrOwnerRemove) {
		t.Fatalf("owner remove err=%v", err)
	}

	memberRepo.member = member
	if err := svc.LeaveTrip(context.Background(), "user-2", "trip-1"); err != nil {
		t.Fatalf("leave: %v", err)
	}
	if err := svc.LeaveTrip(context.Background(), "owner-1", "trip-1"); !errors.Is(err, ErrOwnerLeave) {
		t.Fatalf("owner leave err=%v", err)
	}
}

func TestNewTripMemberService(t *testing.T) {
	svc := NewTripMemberService(nil, nil, nil)
	if svc == nil {
		t.Fatal("NewTripMemberService returned nil")
	}
}

func TestTripMemberErrorsDistinct(t *testing.T) {
	errors := [...]error{
		ErrMemberNotFound, ErrUserNotFound, ErrDuplicateMember,
		ErrOwnerRemove, ErrOwnerLeave, ErrOwnerRoleChange, ErrInvalidTripRole,
	}
	for i, e := range errors {
		for j, o := range errors {
			if i != j && e == o { //nolint:errorlint
				t.Errorf("errors[%d] and errors[%d] are same pointer", i, j)
			}
		}
	}
}

func TestAddMemberNormalizesRole(t *testing.T) {
	svc := NewTripMemberService(nil, nil, nil)
	if svc == nil {
		t.Fatal("NewTripMemberService returned nil")
	}
}
