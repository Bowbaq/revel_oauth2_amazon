package controllers

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/robfig/revel"
	"net/http"
	"net/url"
	"revel_oauth2_amazon/app/models"
	"strconv"
)

type App struct {
	*revel.Controller
}

var AMAZON = &oauth.Config{
	ClientId:     "amzn1.application-oa2-client.b0ca63fd94cd4e0ab836d61529fa49ac",
	ClientSecret: "6755fce240af835e2461c2352480ab713fbc6ee01ee17dedeb88ffdbc1ad281b",
	Scope:        "profile",
	AuthURL:      "https://www.amazon.com/ap/oa",
	TokenURL:     "https://api.amazon.com/auth/o2/token",
	RedirectURL:  "https://revel-oauth2-amazon.herokuapp.com/login",
}

func (app App) Index() revel.Result {
	user := app.currentUser()
	profile := map[string]interface{}{}
	if user != nil && user.Token.AccessToken != "" {
		resp, _ := http.Get("https://api.amazon.com/user/profile?access_token=" +
			url.QueryEscape(user.Token.AccessToken))
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
			revel.ERROR.Println(err)
		}
		revel.INFO.Println(profile)
	}

	authUrl := AMAZON.AuthCodeURL("access_request")
	return app.Render(profile, authUrl)
}

func (app App) Login(code string) revel.Result {
	transport := &oauth.Transport{Config: AMAZON}
	token, err := transport.Exchange(code)
	if err != nil {
		revel.ERROR.Println(err)
		return app.Redirect(App.Index)
	}

	user := app.currentUser()
	models.SetToken(user.Uid, token)

	return app.Redirect(App.Index)
}

func setUser(app *revel.Controller) revel.Result {
	var user *models.User
	if _, ok := app.Session["uid"]; ok {
		uid, _ := strconv.ParseInt(app.Session["uid"], 10, 0)
		user = models.GetUser(uint64(uid))
	}
	if user == nil {
		user = models.NewUser()
		app.Session["uid"] = fmt.Sprintf("%d", user.Uid)
	}
	app.RenderArgs["user"] = user
	return nil
}

func init() {
	revel.InterceptFunc(setUser, revel.BEFORE, &App{})
}

func (c App) currentUser() *models.User {
	return c.RenderArgs["user"].(*models.User)
}
