package auth

import "testing"

func TestPasswordPolicy(t *testing.T) {
	if err := ValidatePassword("short"); err == nil {
		t.Fatal("expected weak password to be rejected")
	}

	if err := ValidatePassword("StrongPass123!"); err != nil {
		t.Fatalf("expected strong password to be accepted: %v", err)
	}
}

func TestPasswordHashRoundTrip(t *testing.T) {
	hasher := NewPasswordHasher()
	hash, err := hasher.HashPassword("StrongPass123!")
	if err != nil {
		t.Fatalf("hashing password failed: %v", err)
	}

	if !hasher.VerifyPassword(hash, "StrongPass123!") {
		t.Fatal("expected password verification to succeed")
	}

	if hasher.VerifyPassword(hash, "WrongPass123!") {
		t.Fatal("expected password verification to fail for wrong password")
	}
}
