package client

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	host       string
}

// UserParam defines GitHub's /user response.
type UserParam struct {
	Login                   string `json:"login"`
	AvatarURL               string `json:"avatar_url"`
	Type                    string `json:"type"` // must to "User"
	Name                    string `json:"name"` // must set
	TwoFactorAuthentication bool   `json:"two_factor_authentication"`
}

// EmailParam defines GitHub's /user/emails response.
type EmailParam struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient, "https://api.github.com"}
}

// User retrieves github user information.
func (c *Client) User() (*UserParam, error) {
	userResponse, err := c.httpClient.Get(c.host + "/user")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := userResponse.Body.Close(); err != nil {
			log.Printf("failed to close github response body: %v", err)
		}
	}()
	var userData UserParam
	if err := json.NewDecoder(userResponse.Body).Decode(&userData); err != nil {
		return nil, err
	}
	return &userData, nil
}

// Emails retrieves list of github user emails.
func (c *Client) Emails() ([]EmailParam, error) {
	emailsResponse, err := c.httpClient.Get(c.host + "/user/emails")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := emailsResponse.Body.Close(); err != nil {
			log.Printf("failed to close github response body: %v", err)
		}
	}()
	var emailsData []EmailParam
	if err := json.NewDecoder(emailsResponse.Body).Decode(&emailsData); err != nil {
		return nil, err
	}
	return emailsData, nil
}

func (u *UserParam) ValidateType() bool {
	return u.Type == "User"
}

func (u *UserParam) ValidateName() bool {
	trim := strings.Trim(u.Name, " ")
	if len(trim) < 3 {
		return false
	}
	if !strings.Contains(trim, " ") {
		return false
	}
	return true
}

func (e *EmailParam) ValidateEmail(keyword string) bool {
	return strings.HasSuffix(e.Email, keyword)
}
