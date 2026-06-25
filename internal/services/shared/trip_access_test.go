package serviceshared

import "testing"

func TestCanViewTrip(t *testing.T) {
	if !CanViewTrip("owner") {
		t.Error("owner should view")
	}
	if !CanViewTrip("editor") {
		t.Error("editor should view")
	}
	if !CanViewTrip("viewer") {
		t.Error("viewer should view")
	}
	if CanViewTrip("") {
		t.Error("empty should not view")
	}
	if CanViewTrip("unknown") {
		t.Error("unknown should not view")
	}
}

func TestCanEditTrip(t *testing.T) {
	if !CanEditTrip("owner") {
		t.Error("owner should edit")
	}
	if !CanEditTrip("editor") {
		t.Error("editor should edit")
	}
	if CanEditTrip("viewer") {
		t.Error("viewer should not edit")
	}
	if CanEditTrip("") {
		t.Error("empty should not edit")
	}
}

func TestCanManageTripMembers(t *testing.T) {
	if !CanManageTripMembers("owner") {
		t.Error("owner should manage")
	}
	if CanManageTripMembers("editor") {
		t.Error("editor should not manage")
	}
	if CanManageTripMembers("viewer") {
		t.Error("viewer should not manage")
	}
}
