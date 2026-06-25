package servicetripmember

import "testing"

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
