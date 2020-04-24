package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mikan/github-id-checker/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const defaultPort = "8080"
const emailMissing = "該当メールアドレスなし"

var org = os.Getenv("ORG")
var keyword = os.Getenv("KEYWORD")
var policyURL = os.Getenv("POLICY_URL")
var webhook = os.Getenv("WEBHOOK")
var config = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
	Scopes:       []string{"user"},
}

// main executes the server program.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/github", handleCallback)
	http.HandleFunc("/submit", handleSubmit)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./template/assets"))))
	log.Printf("Server starts at port %s", port)
	log.Println(http.ListenAndServe(":"+port, nil))
}

// handleIndex handles index page request.
func handleIndex(w http.ResponseWriter, _ *http.Request) {
	t := template.Must(template.New("index.html").ParseFiles("template/index.html"))
	if err := t.Execute(w, struct {
		Org       string
		PolicyURL string
	}{
		org,
		policyURL,
	}); err != nil {
		log.Printf("failed to execute index template: %v", err)
	}
}

// handleLogin handles authentication and authorization request.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL(""), http.StatusTemporaryRedirect)
}

// handleCallback handles OAuth response and transform the information.
func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("failed to exchange code: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	c := client.NewClient(config.Client(context.Background(), token))

	// gather user information
	userData, err := c.User()
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to get user information: %v\n", err))
		return
	}

	// gather user emails
	emailsData, err := c.Emails()
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to get user information: %v\n", err))
		return
	}

	// gather check results
	validEmail := findValidEmail(emailsData, keyword)
	nameChecked := userData.ValidateName()
	typeChecked := userData.ValidateType()
	emailChecked := validEmail.Email != emailMissing
	allChecked := nameChecked && typeChecked && emailChecked && validEmail.Verified && userData.TwoFactorAuthentication
	t := template.Must(template.New("github.html").ParseFiles("template/github.html"))
	if err := t.Execute(w, struct {
		Org          string
		User         *client.UserParam
		Email        *client.EmailParam
		TypeChecked  bool
		NameChecked  bool
		EmailChecked bool
		AllChecked   bool
		AccessToken  string
	}{
		org,
		userData,
		validEmail,
		typeChecked,
		nameChecked,
		emailChecked,
		allChecked,
		token.AccessToken,
	}); err != nil {
		log.Printf("failed to execute github template: %v", err)
	}
}

// handleSubmit handles submit request.
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// retake user information
	c := client.NewClient(config.Client(context.Background(), &oauth2.Token{AccessToken: r.FormValue("token")}))
	userData, err := c.User()
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to get user information: %v\n", err))
		return
	}

	// send webhook notification
	msg := client.Message{
		Text: fmt.Sprintf("%s GitHub ユーザー <https://github.com/%s|%s> から登録依頼が来ました。\n"+
			"<https://github.com/orgs/%s/people|メンバー管理はこちら>",
			os.Getenv("MSG_PREFIX"), userData.Login, userData.Login, org),
	}
	err = msg.Send(webhook)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to send webhook notification: %v", err))
		return
	}

	// send mail
	go sendMail(userData.Login, org)

	t := template.Must(template.New("success.html").ParseFiles("template/success.html"))
	if err := t.Execute(w, struct{ Org string }{Org: org}); err != nil {
		log.Printf("failed to execute success template: %v", err)
	}
}

func findValidEmail(emails []client.EmailParam, keyword string) *client.EmailParam {
	for _, email := range emails {
		if email.ValidateEmail(keyword) {
			return &email
		}
	}
	return &client.EmailParam{Email: emailMissing}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("showing %d error page in %s: %s", status, r.RequestURI, message)
	w.WriteHeader(status)
	t := template.Must(template.New("error.html").ParseFiles("template/error.html"))
	if err := t.Execute(w, struct {
		Org     string
		Status  string
		Message string
	}{
		Org:     org,
		Status:  fmt.Sprintf("HTTP %d %s", status, http.StatusText(status)),
		Message: message,
	}); err != nil {
		log.Printf("failed to execute error template: %v", err)
	}
}
