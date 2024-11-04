package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ResponseBody represents the JSON structure
type ResponseBody struct {
	Id        string `json:"id"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	LastLogin string `json:"lastLogin"`
}

func GetUserEmailById(userId string) (string, error) {

	fmt.Println(userId)
	// Make a GET request to the URL
	resp, err := http.Get(os.Getenv("USER_SERVICE_URL") + "/users/" + userId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	fmt.Println(responseBody.Email, err)

	return responseBody.Email, nil
}
