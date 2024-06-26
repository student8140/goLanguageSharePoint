package sharepointAccess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Replace these values with your own
const (
	tenantID     = "your-tenant-id"
	clientID     = "your-client-id"
	clientSecret = "your-client-secret"
	scope        = "https://your-sharepoint-site.sharepoint.com/.default"
	tokenURL     = "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token"
	sharePointURL = "https://your-sharepoint-site.sharepoint.com/sites/your-site/_api/web/lists/getbytitle('YourListName')/items"
)

// Struct for parsing the token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	// Get OAuth token
	token, err := getOAuthToken()
	if err != nil {
		log.Fatalf("Error getting OAuth token: %v", err)
	}

	// Use token to make an authenticated request to SharePoint
	err = getSharePointListItems(token)
	if err != nil {
		log.Fatalf("Error getting SharePoint list items: %v", err)
	}
}

func getOAuthToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", scope)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	return tokenResponse.AccessToken, nil
}

func getSharePointListItems(token string) error {
	req, err := http.NewRequest("GET", sharePointURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	fmt.Println("SharePoint List Items:", string(body))
	return nil
}
