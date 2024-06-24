package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Azure SharePoint list credentials and endpoint
const (
	sharePointURL = "https://your-sharepoint-site.sharepoint.com/sites/your-site/_api/web/lists/getbytitle('YourListName')/items"
	username      = "your-username@your-domain.com"
	password      = "your-password"
)

// JFrog Artifactory API base URL and API key
const (
	artifactoryBaseURL = "https://your-artifactory-url/artifactory/api"
	artifactoryAPIKey  = "your-api-key"
)

// Struct to unmarshal SharePoint list item JSON
type SharePointListItem struct {
	Title     string `json:"Title"`
	RepoKey   string `json:"RepoKey"`
	GroupName string `json:"GroupName"`
	PermName  string `json:"PermissionName"`
}

// Function to fetch SharePoint list items
func fetchSharePointListItems() ([]SharePointListItem, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", sharePointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating SharePoint request: %v", err)
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing SharePoint request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading SharePoint response body: %v", err)
	}

	var result struct {
		Value []SharePointListItem `json:"value"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing SharePoint response: %v", err)
	}

	return result.Value, nil
}

// CreateRepository creates a new repository in JFrog Artifactory
func CreateRepository(repoKey string) error {
	repoBody := map[string]interface{}{
		"rclass":      "local",
		"packageType": "generic",
		"key":         repoKey,
	}

	body, _ := json.Marshal(repoBody)
	req, err := http.NewRequest("PUT", artifactoryBaseURL+"/repositories/"+repoKey, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating repository request: %v", err)
	}

	req.Header.Set("X-JFrog-Art-Api", artifactoryAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating repository: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Repository created:", resp.Status, string(respBody))

	return nil
}

// CreateGroup creates a new group in JFrog Artifactory
func CreateGroup(groupName string) error {
	groupBody := map[string]interface{}{
		"name":        groupName,
		"description": "This is a new group",
	}

	body, _ := json.Marshal(groupBody)
	req, err := http.NewRequest("PUT", artifactoryBaseURL+"/security/groups/"+groupName, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating group request: %v", err)
	}

	req.Header.Set("X-JFrog-Art-Api", artifactoryAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Group created:", resp.Status, string(respBody))

	return nil
}

// CreatePermissionTarget creates a new permission target in JFrog Artifactory
func CreatePermissionTarget(permissionName, repoKey, groupName string) error {
	permissionBody := map[string]interface{}{
		"name": permissionName,
		"repositories": []string{
			repoKey,
		},
		"principals": map[string]interface{}{
			"groups": map[string][]string{
				groupName: {"r", "w", "m", "d", "n"},
			},
		},
	}

	body, _ := json.Marshal(permissionBody)
	req, err := http.NewRequest("PUT", artifactoryBaseURL+"/security/permissions/"+permissionName, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating permission target request: %v", err)
	}

	req.Header.Set("X-JFrog-Art-Api", artifactoryAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error creating permission target: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Permission target created:", resp.Status, string(respBody))

	return nil
}

func main() {
	// Fetch SharePoint list items
	items, err := fetchSharePointListItems()
	if err != nil {
		log.Fatalf("Error fetching SharePoint list items: %v", err)
	}

	// Process each item
	for _, item := range items {
		err := CreateRepository(item.RepoKey)
		if err != nil {
			log.Printf("Error creating repository %s: %v", item.RepoKey, err)
		}

		err = CreateGroup(item.GroupName)
		if err != nil {
			log.Printf("Error creating group %s: %v", item.GroupName, err)
		}

		err = CreatePermissionTarget(item.PermName, item.RepoKey, item.GroupName)
		if err != nil {
			log.Printf("Error creating permission target %s: %v", item.PermName, err)
		}
	}
}
