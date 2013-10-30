package middleware

import (
    "fmt"
    //"log"
    "errors"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/stevan/httpapp"
    "github.com/stevan/httpapp/middleware/sessions"

    "code.google.com/p/goauth2/oauth"
)

// Simple user objects ...

type GoogleUserInfo struct {
    Id       string `json:"id"`
    Email    string `json:"email"`
    Verified bool   `json:"verified_email"`
}

// ... Auth middleware

type GoogleOAuthHandler struct {
    App httpapp.App
    OauthConfig *oauth.Config
}

func (a *GoogleOAuthHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    a.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (a *GoogleOAuthHandler) Call (e *httpapp.Env) *httpapp.Response {
    path := e.Request.URL.Path
    switch path {
        case "/oauth2callback":
            return a.handleOAuth2Callback(e)
        default:
            return a.checkSession(e)
    }
}

func (a *GoogleOAuthHandler) checkSession (e *httpapp.Env) *httpapp.Response {
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

func (a *GoogleOAuthHandler) handleOAuth2Callback (e *httpapp.Env) *httpapp.Response {

    t := &oauth.Transport{ Config: a.OauthConfig }
    t.Exchange(e.Request.FormValue("code"))

    req, err := t.Client().Get("https://www.googleapis.com/oauth2/v1/userinfo")
    if err != nil { panic(err) }
    defer req.Body.Close()

    body, _ := ioutil.ReadAll(req.Body)

    //log.Println(string(body))

    var user_info GoogleUserInfo
    json_err := json.Unmarshal(body, &user_info)
    if json_err != nil { panic(json_err) }

    if user_info.Email == "" {
        panic(errors.New(fmt.Sprintf("Incorrect JSON response from Google: %s", body)))
    }

    //log.Printf("%v", user_info)

    session := e.Get("session").(*sessions.Session)
    session.Data["user"] = &user_info

    resp := httpapp.NewResponse(http.StatusFound)
    resp.Headers.Add("Location", session.Data["return_to"].(string))
    return resp
}







