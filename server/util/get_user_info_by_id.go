package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

func GetUserInfoById(userId string) (ResponseBody, error) {

	var responseBody ResponseBody

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Make a GET request to the URL
	resp, err := http.Get(os.Getenv("USER_SERVICE_URL") + "/users/" + userId)
	if err != nil {
		return responseBody, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return responseBody, fmt.Errorf("response status code is not ok; received: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseBody, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return responseBody, fmt.Errorf("failed to parse JSON: %v", err)
	}

	fmt.Println(responseBody, err)

	return responseBody, nil
}
