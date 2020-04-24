package main

import (
	"testing"

	"github.com/mikan/github-id-checker/client"
)

func TestFindValidEmail(t *testing.T) {
	emails := []client.EmailParam{{"test1@test.com", false}, {"test2@abc.com", true}}
	valid := findValidEmail(emails, "abc.com")
	if valid.Email != "test2@abc.com" {
		t.Fatalf("Email: expected %s, actual %s", "test2@abc.com", valid.Email)
	}
	if !valid.Verified {
		t.Fatalf("Verified: expected %v, actual %v", true, valid.Verified)
	}
}
