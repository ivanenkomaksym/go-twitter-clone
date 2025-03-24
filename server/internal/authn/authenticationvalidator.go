package authn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"
)

type IAuthenticationValidator interface {
	ValidateAuthentication(w http.ResponseWriter, r *http.Request) *models.User
}

type AuthenticationValidator struct {
	Authentication config.Authentication
}

func (validator AuthenticationValidator) ValidateAuthentication(w http.ResponseWriter, r *http.Request) *models.User {
	if !validator.Authentication.Enable {
		return &models.User{IsAnonymous: true}
	}

	id_token := r.Header.Get("Authorization")
	if len(id_token) > 7 && id_token[:7] == "Bearer " {
		id_token = id_token[7:]
	}
	if id_token != "" {
		url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", id_token)
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Failed to validate id_token", http.StatusUnauthorized)
			return nil
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response from Google", http.StatusUnauthorized)
			return nil
		}

		// Parse the response body (which contains user info)
		var response map[string]any
		if err := json.Unmarshal(body, &response); err != nil {
			http.Error(w, "Failed to parse user info", http.StatusUnauthorized)
			return nil
		}

		user := models.User{
			IsAnonymous: false,
			FirstName:   response["given_name"].(string),
			LastName:    response["family_name"].(string),
			Email:       response["email"].(string),
			Picture:     response["picture"].(string),
		}

		return &user
	}

	cookie, err := r.Cookie("id_token")
	if err != nil {
		http.Error(w, "Unauthorized: No id_token found", http.StatusUnauthorized)
		return nil
	}

	token := cookie.Value
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to validate id_token", http.StatusUnauthorized)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Google", http.StatusUnauthorized)
		return nil
	}

	// Parse the response body (which contains user info)
	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusUnauthorized)
		return nil
	}

	errordesc, haserror := response["error_description"]
	if haserror {
		http.Error(w, errordesc.(string), http.StatusUnauthorized)
		return nil
	}

	user := models.User{
		IsAnonymous: false,
		FirstName:   getString(response["given_name"]),
		LastName:    getString(response["family_name"]),
		Email:       getString(response["email"]),
		Picture:     getString(response["picture"]),
	}

	return &user
}

func getString(value any) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
