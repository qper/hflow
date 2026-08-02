package domain

import "testing"

func TestUserDefaults(t *testing.T) {
	user := User{Username: "alice", Email: "alice@example.com"}

	if user.Username != "alice" {
		t.Fatalf("expected username to be preserved, got %q", user.Username)
	}

	if user.Email != "alice@example.com" {
		t.Fatalf("expected email to be preserved, got %q", user.Email)
	}
}
