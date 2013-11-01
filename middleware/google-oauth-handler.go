package middleware

import (
	"fmt"
	//"log"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/stevan/httpapp"
	"github.com/stevan/httpapp/middleware/auth"
	"github.com/stevan/httpapp/middleware/sessions"

	"code.google.com/p/goauth2/oauth"
)

// ... Auth middleware

type GoogleOAuthHandler struct {
	App         httpapp.App
	OauthConfig *oauth.Config
}

func (a *GoogleOAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (a *GoogleOAuthHandler) Call(e *httpapp.Env) *httpapp.Response {
	path := e.Request.URL.Path
	switch path {
	case "/oauth2callback":
		return a.handleOAuth2Callback(e)
	default:
		return a.checkSession(e)
	}
}

func (a *GoogleOAuthHandler) checkSession(e *httpapp.Env) *httpapp.Response {
	session := e.Get("session").(*sessions.Session)
	if _, ok := session.Data["user"]; !ok {
		session.Data["return_to"] = e.Request.URL.Path
		url := a.OauthConfig.AuthCodeURL("")
		resp := httpapp.NewResponse(http.StatusFound)
		resp.Headers.Add("Location", url)
		return resp
	} else {
		return a.App.Call(e)
	}
}

func (a *GoogleOAuthHandler) handleOAuth2Callback(e *httpapp.Env) *httpapp.Response {

	t := &oauth.Transport{Config: a.OauthConfig}
	t.Exchange(e.Request.FormValue("code"))

	base_user_info := a.getBaseUserInfo(t)
	ext_user_info := a.getExtendedUserInfo(t, base_user_info)
	user_info := auth.NewGoogleUser(
		base_user_info.Id,
		base_user_info.Email,
		ext_user_info.FullName,
	)

	session := e.Get("session").(*sessions.Session)
	session.Data["user"] = user_info

	resp := httpapp.NewResponse(http.StatusFound)
	resp.Headers.Add("Location", session.Data["return_to"].(string))
	return resp
}

// structs for JSON UnMarshal

type googleUserInfoJSON struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
}

type googlePeopleInfoJSON struct {
	Id       string `json:"id"`
	FullName string `json:"displayName"`
	Verified bool   `json:"verified"`
}

func (a *GoogleOAuthHandler) getBaseUserInfo(t *oauth.Transport) *googleUserInfoJSON {
	req, err := t.Client().Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	body, _ := ioutil.ReadAll(req.Body)

	var user_info googleUserInfoJSON
	json_err := json.Unmarshal(body, &user_info)
	if json_err != nil {
		panic(json_err)
	}

	if user_info.Id == "" {
		panic(errors.New(fmt.Sprintf("Incorrect JSON response from Google: %s", body)))
	}

	return &user_info
}

func (a *GoogleOAuthHandler) getExtendedUserInfo(t *oauth.Transport, info *googleUserInfoJSON) *googlePeopleInfoJSON {
	req, err := t.Client().Get(fmt.Sprintf("https://www.googleapis.com/plus/v1/people/%s", info.Id))
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	body, _ := ioutil.ReadAll(req.Body)

	var user_info googlePeopleInfoJSON
	json_err := json.Unmarshal(body, &user_info)
	if json_err != nil {
		panic(json_err)
	}

	if user_info.Id != info.Id {
		panic(errors.New(fmt.Sprintf("Incorrect JSON response from Google: %s", body)))
	}

	return &user_info
}
