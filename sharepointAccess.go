package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Function to get access token
func getAccessToken(clientID, clientSecret, tenantID string) (string, error) {
	authURL := "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "https://your-tenant.sharepoint.com/.default") // Replace 'your-tenant' with your actual tenant name

	req, err := http.NewRequest("POST", authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get access token: %s, response: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return accessToken, nil
}

// Function to get SharePoint list items
func getSharePointList(accessToken, siteURL, listName string) (string, error) {
	apiURL := fmt.Sprintf("%s/_api/web/lists/GetByTitle('%s')/items", siteURL, listName)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept", "application/json;odata=verbose")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get list items: %s, response: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

// Main function
func main() {
	clientID := "your-client-id"
	clientSecret := "your-client-secret"
	tenantID := "your-tenant-id"
	siteURL := "https://your-tenant.sharepoint.com/sites/yoursite" // Replace 'your-tenant' with your actual tenant name
	listName := "your-list-name"

	// Get access token
	token, err := getAccessToken(clientID, clientSecret, tenantID)
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}
	fmt.Println("Access Token:", token)

	// Get SharePoint list items
	listItems, err := getSharePointList(token, siteURL, listName)
	if err != nil {
		log.Fatalf("Error getting SharePoint list: %v", err)
	}

	fmt.Println("SharePoint List Items:", listItems)
}
