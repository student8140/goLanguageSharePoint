package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// RepositoryConfig represents the configuration for a new repository
type RepositoryConfig struct {
	RClass      string `json:"rclass"`
	PackageType string `json:"packageType"`
}

// PermissionTargetConfig represents the configuration for a permission target
type PermissionTargetConfig struct {
	Name         string            `json:"name"`
	Repositories []string          `json:"repositories"`
	Principals   struct {
		Users  map[string][]string `json:"users"`
		Groups map[string][]string `json:"groups"`
	} `json:"principals"`
}

func main() {
	jfrogURL := os.Getenv("JFROG_URL")        // JFrog URL
	jfrogAPIKey := os.Getenv("JFROG_API_KEY") // JFrog API Key

	if jfrogURL == "" || jfrogAPIKey == "" {
		fmt.Println("Please set the JFROG_URL and JFROG_API_KEY environment variables.")
		return
	}

	// Create Repository
	repoConfig := RepositoryConfig{
		RClass:      "local",
		PackageType: "maven",
	}

	repoName := "example-repo-local"
	createRepo(jfrogURL, jfrogAPIKey, repoName, repoConfig)

	// Update Permission Target
	permissionTargetName := "existing-permission"
	userName := "new-user" // User to add
	permissions := []string{"read", "write", "annotate"} // Permissions for the user

	updatePermissionTarget(jfrogURL, jfrogAPIKey, permissionTargetName, repoName, userName, permissions)
}

func createRepo(jfrogURL, jfrogAPIKey, repoName string, config RepositoryConfig) {
	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		return
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/artifactory/api/repositories/%s", jfrogURL, repoName), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %s\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-JFrog-Art-Api", jfrogAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading HTTP response: %s\n", err)
		return
	}

	fmt.Printf("Create Repository - Response Status: %s\n", resp.Status)
	fmt.Printf("Create Repository - Response Body: %s\n", string(body))
}

func updatePermissionTarget(jfrogURL, jfrogAPIKey, permissionTargetName, repoName, userName string, permissions []string) {
	// Get existing permission target
	permissionTarget, err := getPermissionTarget(jfrogURL, jfrogAPIKey, permissionTargetName)
	if err != nil {
		fmt.Printf("Error getting permission target: %s\n", err)
		return
	}

	// Add new repository to the list of repositories
	permissionTarget.Repositories = append(permissionTarget.Repositories, repoName)

	// Add user to the principals if not already present
	if permissionTarget.Principals.Users == nil {
		permissionTarget.Principals.Users = make(map[string][]string)
	}
	permissionTarget.Principals.Users[userName] = permissions

	jsonData, err := json.Marshal(permissionTarget)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		return
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/artifactory/api/security/permissions/%s", jfrogURL, permissionTargetName), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %s\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-JFrog-Art-Api", jfrogAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading HTTP response: %s\n", err)
		return
	}

	fmt.Printf("Update Permission Target - Response Status: %s\n", resp.Status)
	fmt.Printf("Update Permission Target - Response Body: %s\n", string(body))
}

func getPermissionTarget(jfrogURL, jfrogAPIKey, permissionTargetName string) (*PermissionTargetConfig, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/artifactory/api/security/permissions/%s", jfrogURL, permissionTargetName), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("X-JFrog-Art-Api", jfrogAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get permission target. Status: %s, Response: %s", resp.Status, string(body))
	}

	var permissionTarget PermissionTargetConfig
	err = json.Unmarshal(body, &permissionTarget)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &permissionTarget, nil
}
