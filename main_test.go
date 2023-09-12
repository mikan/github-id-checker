package main

import (
	"testing"

	"github.com/mikan/github-id-checker/client"
)

func TestFindValidEmail(t *testing.T) {
	emails := []client.EmailParam{{Email: "test1@test.com", Verified: false}, {Email: "test2@abc.com", Verified: true}}
	valid := findValidEmail(emails, "abc.com")
	if valid.Email != "test2@abc.com" {
		t.Fatalf("Email: expected %s, actual %s", "test2@abc.com", valid.Email)
	}
	if !valid.Verified {
		t.Fatalf("Verified: expected %v, actual %v", true, valid.Verified)
	}
}
