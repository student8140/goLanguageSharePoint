package main

import (
	"fmt"
	"os"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	jfrogURL := os.Getenv("JFROG_URL")
	jfrogAPIKey := os.Getenv("JFROG_API_KEY")

	if jfrogURL == "" || jfrogAPIKey == "" {
		fmt.Println("Please set the JFROG_URL and JFROG_API_KEY environment variables.")
		return
	}

	// Set up Artifactory details
	serviceDetails := auth.NewArtifactoryDetails()
	serviceDetails.SetUrl(jfrogURL)
	serviceDetails.SetApiKey(jfrogAPIKey)

	// Create Artifactory service manager
	log.SetLogger(log.NewLogger(log.INFO, nil))
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(serviceDetails).
		Build()

	if err != nil {
		log.Error(err)
		return
	}

	artifactoryServiceManager, err := artifactory.New(serviceConfig)
	if err != nil {
		log.Error(err)
		return
	}

	// Create repository
	repoName := "example-repo-local"
	createRepo(artifactoryServiceManager, repoName)

	// Update permission target
	permissionTargetName := "existing-permission"
	users := map[string][]string{
		"new-user": {"read", "write", "annotate"},
	}
	groups := map[string][]string{
		"existing-group": {"read", "write"},
	}

	updatePermissionTarget(artifactoryServiceManager, permissionTargetName, repoName, users, groups)

	// Remove group from permission target
	groupName := "group-to-remove"
	removeGroupFromPermissionTarget(artifactoryServiceManager, permissionTargetName, groupName)
}

func createRepo(serviceManager *artifactory.ArtifactoryServicesManager, repoName string) {
	repoConfig := services.LocalRepositoryBaseParams{
		RepositoryBaseParams: services.RepositoryBaseParams{
			Key: repoName,
		},
		PackageType: "maven",
	}

	err := serviceManager.CreateLocalRepository(repoConfig)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Repository %s created successfully\n", repoName)
}

func updatePermissionTarget(serviceManager *artifactory.ArtifactoryServicesManager, permissionTargetName, repoName string, users, groups map[string][]string) {
	permissionTargetParams := services.PermissionTargetParams{
		Name: permissionTargetName,
	}

	permissionTarget, err := serviceManager.GetPermissionTarget(permissionTargetName)
	if err != nil {
		log.Error(err)
		return
	}

	// Add repository to the permission target
	permissionTarget.Repositories = append(permissionTarget.Repositories, repoName)

	// Add users to the permission target
	if permissionTarget.Principals.Users == nil {
		permissionTarget.Principals.Users = make(map[string][]string)
	}
	for user, perms := range users {
		permissionTarget.Principals.Users[user] = perms
	}

	// Add groups to the permission target
	if permissionTarget.Principals.Groups == nil {
		permissionTarget.Principals.Groups = make(map[string][]string)
	}
	for group, perms := range groups {
		permissionTarget.Principals.Groups[group] = perms
	}

	err = serviceManager.UpdatePermissionTarget(*permissionTarget)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Permission target %s updated successfully\n", permissionTargetName)
}

func removeGroupFromPermissionTarget(serviceManager *artifactory.ArtifactoryServicesManager, permissionTargetName, groupName string) {
	permissionTarget, err := serviceManager.GetPermissionTarget(permissionTargetName)
	if err != nil {
		log.Error(err)
		return
	}

	// Remove group from permission target
	delete(permissionTarget.Principals.Groups, groupName)

	err = serviceManager.UpdatePermissionTarget(*permissionTarget)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Group %s removed from permission target %s successfully\n", groupName, permissionTargetName)
}
