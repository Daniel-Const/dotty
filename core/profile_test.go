package core

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	p := NewProfile("./profiles/potato")

	// Should set Location correctly
	if p.Location != "./profiles/potato" {
		t.Fatalf("Profile.Location is %s | Expecting %s", p.Location, "./profiles/potato")
	}

	// Should set Name correctly
	if p.Name != "potato" {
		t.Fatalf("Profile.Name is %s | Expecting %s", p.Name, "potato")
	}
}
