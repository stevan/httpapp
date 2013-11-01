package auth

import (
	"fmt"

	"code.google.com/p/goauth2/oauth"
)

// -----------------------------------------
// Config Object
// -----------------------------------------

func CreateGoogleOAuthConfig(client_id string, client_secret string, redirect_URL string, cache_file_path string) *oauth.Config {
	return &oauth.Config{
		ClientId:     client_id,
		ClientSecret: client_secret,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		RedirectURL:  redirect_URL,
		Scope:        "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/plus.me",
		TokenCache:   oauth.CacheFile(cache_file_path),
	}
}

// -----------------------------------------
// User Object
// -----------------------------------------

func NewGoogleUser(id string, email string, full_name string) *GoogleUser {
	return &GoogleUser{id: id, email: email, full_name: full_name}
}

// User Object

type GoogleUser struct {
	id        string
	email     string
	full_name string
}

func (u *GoogleUser) Email() string {
	return u.email
}

func (u *GoogleUser) FullName() string {
	return u.full_name
}

func (u *GoogleUser) AuthorString() string {
	return fmt.Sprintf("%s <%s>", u.full_name, u.email)
}
