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

	// Create new group
	groupName := "new-group"
	createGroup(artifactoryServiceManager, groupName)

	// Create new permission target
	permissionTargetName := "new-permission"
	users := map[string][]string{
		"new-user": {"read", "write", "annotate"},
	}
	groups := map[string][]string{
		groupName: {"read", "write"},
	}

	createPermissionTarget(artifactoryServiceManager, permissionTargetName, repoName, users, groups)

	// Add group to existing permission target
	existingPermissionTargetName := "existing-permission"
	addGroupToPermissionTarget(artifactoryServiceManager, existingPermissionTargetName, groupName, []string{"read", "write"})
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

func createGroup(serviceManager *artifactory.ArtifactoryServicesManager, groupName string) {
	groupParams := services.GroupParams{
		Name: groupName,
	}

	err := serviceManager.CreateGroup(groupParams)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Group %s created successfully\n", groupName)
}

func createPermissionTarget(serviceManager *artifactory.ArtifactoryServicesManager, permissionTargetName, repoName string, users, groups map[string][]string) {
	permissionTargetParams := services.PermissionTargetParams{
		Name: permissionTargetName,
		Repositories: []string{
			repoName,
		},
		Principals: services.Principals{
			Users:  users,
			Groups: groups,
		},
	}

	err := serviceManager.CreatePermissionTarget(permissionTargetParams)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Permission target %s created successfully\n", permissionTargetName)
}

func addGroupToPermissionTarget(serviceManager *artifactory.ArtifactoryServicesManager, permissionTargetName, groupName string, permissions []string) {
	permissionTarget, err := serviceManager.GetPermissionTarget(permissionTargetName)
	if err != nil {
		log.Error(err)
		return
	}

	// Add group to the permission target
	if permissionTarget.Principals.Groups == nil {
		permissionTarget.Principals.Groups = make(map[string][]string)
	}
	permissionTarget.Principals.Groups[groupName] = permissions

	err = serviceManager.UpdatePermissionTarget(*permissionTarget)
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Group %s added to permission target %s successfully\n", groupName, permissionTargetName)
}
