package authn

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"twitter-clone/internal/config"

	"golang.org/x/oauth2"
)

type OAuth2Router struct {
	Config config.Configuration
}

func (router OAuth2Router) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := router.generateStateOauthCookie(w)
	u := router.Config.OAuth2.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (router OAuth2Router) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	data, err := router.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func (router OAuth2Router) getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := router.Config.OAuth2.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	log.Println(contents)
	log.Println(token)

	return contents, nil
}

func (router OAuth2Router) generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}