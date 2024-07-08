package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	tenantID          = "your-tenant-id"
	clientID          = "your-client-id"
	clientSecret      = "your-client-secret"
	sharePointSiteURL = "https://your-sharepoint-site-url"
	jFrogURL          = "https://your-jfrog-url"
	jFrogAPIToken     = "your-jfrog-api-token"
)

type SharePointItem struct {
	PackageType     string `json:"PackageType"`
	RepoKey         string `json:"RepoKey"`
	RepoDescription string `json:"RepoDescription"`
	GroupName       string `json:"GroupName"`
	GroupDescription string `json:"GroupDescription"`
	Realm           string `json:"Realm"`
	PermissionName  string `json:"PermissionName"`
}

func getOAuth2Token() string {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("resource", "https://your-sharepoint-site-url")

	req, err := http.NewRequest("POST", fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenantID), bytes.NewBufferString(form.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from OAuth2: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return tokenResponse["access_token"].(string)
}

func getSharePointList(accessToken string) []SharePointItem {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/_api/web/lists/getbytitle('YourListName')/items", sharePointSiteURL), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from SharePoint: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var result struct {
		D struct {
			Results []SharePointItem `json:"results"`
		} `json:"d"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return result.D.Results
}

func createJFrogRepository(repoKey, packageType, repoDescription string) {
	client := &http.Client{}

	repoData := map[string]interface{}{
		"rclass":       "local",
		"packageType":  packageType,
		"repoKey":      repoKey,
		"description":  repoDescription,
	}
	repoJSON, err := json.Marshal(repoData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/artifactory/api/repositories/%s", jFrogURL, repoKey), bytes.NewBuffer(repoJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jFrogAPIToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from JFrog: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Printf("JFrog repository creation response: %s\n", string(body))
}

func createJFrogGroup(groupName, groupDescription, realm string) {
	client := &http.Client{}

	groupData := map[string]interface{}{
		"name":        groupName,
		"description": groupDescription,
		"realm":       realm,
	}
	groupJSON, err := json.Marshal(groupData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/artifactory/api/security/groups/%s", jFrogURL, groupName), bytes.NewBuffer(groupJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jFrogAPIToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from JFrog: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Printf("JFrog group creation response: %s\n", string(body))
}

func createJFrogPermission(permissionName, repoKey, groupName string) {
	client := &http.Client{}

	permissionData := map[string]interface{}{
		"name": permissionName,
		"repo": map[string]interface{}{
			"include-patterns": []string{"**"},
			"exclude-patterns": []string{""},
			"repositories":     []string{repoKey},
		},
		"principals": map[string]interface{}{
			"users":  map[string]interface{}{},
			"groups": map[string]interface{}{groupName: []string{"r", "w", "n"}},
		},
	}
	permissionJSON, err := json.Marshal(permissionData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/artifactory/api/security/permissions/%s", jFrogURL, permissionName), bytes.NewBuffer(permissionJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jFrogAPIToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from JFrog: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Printf("JFrog permission creation response: %s\n", string(body))
}

func main() {
	// Get OAuth2 token for SharePoint
	accessToken := getOAuth2Token()

	// Access SharePoint list and get details
	sharePointItems := getSharePointList(accessToken)

	for _, item := range sharePointItems {
		// Create JFrog repository
		createJFrogRepository(item.RepoKey, item.PackageType, item.RepoDescription)

		// Create JFrog group
		createJFrogGroup(item.GroupName, item.GroupDescription, item.Realm)

		// Create JFrog permission
		createJFrogPermission(item.PermissionName, item.RepoKey, item.GroupName)
	}
}
