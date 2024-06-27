package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get access token: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result["access_token"].(string), nil
}

// Function to get SharePoint list items
func getSharePointList(accessToken, siteURL, listName string) (string, error) {
	apiURL := fmt.Sprintf("%s/_api/web/lists/GetByTitle('%s')/items", siteURL, listName)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept", "application/json;odata=verbose")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get list items: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
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
		fmt.Println("Error getting access token:", err)
		return
	}

	// Get SharePoint list items
	listItems, err := getSharePointList(token, siteURL, listName)
	if err != nil {
		fmt.Println("Error getting SharePoint list:", err)
		return
	}

	fmt.Println("SharePoint List Items:", listItems)
}
