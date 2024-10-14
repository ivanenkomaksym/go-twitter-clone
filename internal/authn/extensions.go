package authn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"twitter-clone/internal/models"
)

func ValidateAuthentication(w http.ResponseWriter, r *http.Request) *models.User {
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Google", http.StatusUnauthorized)
		return nil
	}

	// Parse the response body (which contains user info)
	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusUnauthorized)
		return nil
	}

	user := models.User{
		FirstName: userInfo["given_name"].(string),
		LastName:  userInfo["family_name"].(string),
		Email:     userInfo["email"].(string),
		Picture:   userInfo["picture"].(string),
	}

	return &user
}
