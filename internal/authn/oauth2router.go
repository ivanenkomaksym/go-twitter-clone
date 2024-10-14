package authn

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"twitter-clone/internal/config"
	"twitter-clone/internal/extensions"
	"twitter-clone/internal/models"

	"golang.org/x/oauth2"
)

type OAuth2Router struct {
	Config config.Configuration
}

type AuthenticationResult struct {
	Contents []byte
	IdToken  string
}

func (router OAuth2Router) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := router.generateStateOauthCookie(w)
	u := router.Config.OAuth2.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (router OAuth2Router) OauthGoogleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "id_token",
		Value:  "",
		MaxAge: -1, // This deletes the cookie
	})

	// Redirect to the frontend or some protected page
	http.Redirect(w, r, router.Config.RedirectURI, http.StatusFound)
}

func (router OAuth2Router) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	data, err := router.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Set the token in an HTTP-only, secure cookie
	// Additionally token could be encrypted with AES encryption
	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    data.IdToken,
		HttpOnly: true,                 // Prevent JavaScript access
		Secure:   true,                 // Ensure it's sent only over HTTPS
		SameSite: http.SameSiteLaxMode, // Helps mitigate CSRF
	})

	// Redirect to the frontend or some protected page
	http.Redirect(w, r, router.Config.RedirectURI, http.StatusFound)

	fmt.Fprintf(w, "UserInfo: %s\n", data.Contents)
}

func (router OAuth2Router) OauthUserInfo(w http.ResponseWriter, r *http.Request) {
	extensions.EnableCors(&w, router.Config)
	// Get the id_token from the HttpOnly cookie
	cookie, err := r.Cookie("id_token")
	if err != nil {
		http.Error(w, "Unauthorized: No id_token found", http.StatusUnauthorized)
		return
	}

	// Prepare the request to Google's tokeninfo endpoint
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", cookie.Value)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to validate id_token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Google", http.StatusInternalServerError)
		return
	}

	// Parse the response body (which contains user info)
	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	user := models.User{
		FirstName: userInfo["given_name"].(string),
		LastName:  userInfo["family_name"].(string),
		Email:     userInfo["email"].(string),
		Picture:   userInfo["picture"].(string),
	}

	// Return the User info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (router OAuth2Router) getUserDataFromGoogle(code string) (*AuthenticationResult, error) {
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

	log.Println("contents: ", contents)
	log.Println("token: ", token)

	result := AuthenticationResult{
		Contents: contents,
		IdToken:  token.Extra("id_token").(string),
	}

	return &result, nil
}

func (router OAuth2Router) generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}
