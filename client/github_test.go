package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser(t *testing.T) {
	// launch test server
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("must use %s method", http.MethodGet)
			}
			if _, err := fmt.Fprint(w, sampleUserResponse); err != nil {
				t.Logf("failed to write response: %v", err)
			}
		},
	))
	defer ts.Close()

	// call
	c := NewClient(http.DefaultClient)
	c.host = ts.URL
	userData, err := c.User()
	if err != nil {
		t.Fatal(err)
	}
	if !userData.ValidateType() {
		t.Fatalf("Type: invalid")
	}
	if !userData.ValidateName() {
		t.Fatalf("Name: invalid")
	}
	if userData.Name != "monalisa octocat" {
		t.Fatalf("Name: expected %s, actual %s", "monalisa octocat", userData.Name)
	}
}

func TestEmails(t *testing.T) {
	// launch test server
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("must use %s method", http.MethodGet)
			}
			if _, err := fmt.Fprint(w, sampleEmailsResponse); err != nil {
				t.Logf("failed to write response: %v", err)
			}
		},
	))
	defer ts.Close()

	// call
	c := NewClient(http.DefaultClient)
	c.host = ts.URL
	emailsData, err := c.Emails()
	if err != nil {
		t.Fatal(err)
	}
	if emailsData == nil {
		t.Fatal("Emails: nil")
	}
	if len(emailsData) != 1 {
		t.Fatalf("Emails: expected %d, actual %d", 1, len(emailsData))
	}
	emailData := emailsData[0]
	if !emailData.ValidateEmail("github.com") {
		t.Fatal("Emails: not contains github.com")
	}
}

func TestValidateName_validCase(t *testing.T) {
	invalidSamples := []UserParam{
		{Name: "Test User"},
		{Name: "test user"},
		{Name: "T U"},
	}
	for _, sample := range invalidSamples {
		if !sample.ValidateName() {
			t.Fatalf("ValidateName: valid sample rejected: %s", sample.Name)
		}
	}
}

func TestValidateName_invalidCase(t *testing.T) {
	invalidSamples := []UserParam{
		{Name: ""},
		{Name: "1"},
		{Name: "2="},
		{Name: "3=="},
		{Name: "test.user"},
	}
	for _, sample := range invalidSamples {
		if sample.ValidateName() {
			t.Fatalf("ValidateName: invalid sample accepted: %s", sample.Name)
		}
	}
}

const sampleUserResponse = `
{
  "login": "octocat",
  "id": 1,
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/octocat",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/octocat/followers",
  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
  "organizations_url": "https://api.github.com/users/octocat/orgs",
  "repos_url": "https://api.github.com/users/octocat/repos",
  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
  "received_events_url": "https://api.github.com/users/octocat/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa octocat",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "octocat@github.com",
  "hireable": false,
  "bio": "There once was...",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z"
}
`

const sampleEmailsResponse = `
[
  {
    "email": "octocat@github.com",
    "verified": true,
    "primary": true,
    "visibility": "public"
  }
]
`
