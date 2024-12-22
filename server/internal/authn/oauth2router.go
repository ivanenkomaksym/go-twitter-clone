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

	"golang.org/x/oauth2"
)

type OAuth2Router struct {
	Authentication          config.Authentication
	RedirectURI             string
	AllowOrigin             string
	Domain                  string
	AuthenticationValidator IAuthenticationValidator
}

func enableCors(w *http.ResponseWriter, allowOrigin string) {
	(*w).Header().Set("Access-Control-Allow-Origin", allowOrigin)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

type AuthenticationResult struct {
	Contents []byte
	IdToken  string
}

func (router OAuth2Router) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, router.AllowOrigin)

	oauthState := router.generateStateOauthCookie(w)
	u := router.Authentication.OAuth2.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (router OAuth2Router) OauthGoogleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "id_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // This deletes the cookie
	})

	// Redirect to the frontend or some protected page
	http.Redirect(w, r, router.RedirectURI, http.StatusFound)
}

func (router OAuth2Router) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, router.AllowOrigin)

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
		Path:     "/",
		HttpOnly: true, // Prevent JavaScript access
		Secure:   true, // Ensure it's sent only over HTTPS
		SameSite: http.SameSiteNoneMode,
		Domain:   router.Domain,
	})

	fmt.Printf("Set-Cookie: %s\n", w.Header().Get("Set-Cookie"))

	// Redirect to the frontend or some protected page
	http.Redirect(w, r, router.RedirectURI, http.StatusFound)

	fmt.Fprintf(w, "UserInfo: %s\n", data.Contents)
}

func (router OAuth2Router) OauthUserInfo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, router.AllowOrigin)

	user := router.AuthenticationValidator.ValidateAuthentication(w, r)
	if user == nil {
		return
	}

	// Return the User info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (router OAuth2Router) getUserDataFromGoogle(code string) (*AuthenticationResult, error) {
	// Use code to get token and get user info from Google.
	token, err := router.Authentication.OAuth2.Exchange(context.Background(), code)
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
