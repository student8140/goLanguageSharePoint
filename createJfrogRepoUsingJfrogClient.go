package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	// Retrieve the JFrog Artifactory URL and API key from environment variables
	artifactoryURL := os.Getenv("JFROG_ARTIFACTORY_URL")
	apiKey := os.Getenv("JFROG_API_KEY")

	if artifactoryURL == "" || apiKey == "" {
		log.Fatal("Environment variables JFROG_ARTIFACTORY_URL and JFROG_API_KEY must be set")
	}

	// Configure the Artifactory client with API key
	artifactoryDetails := artifactory.NewArtifactoryDetails()
	artifactoryDetails.SetUrl(artifactoryURL)
	artifactoryDetails.SetApiKey(apiKey)

	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(artifactoryDetails).
		SetLog(log.NewLogger(log.INFO, nil)).
		Build()
	if err != nil {
		log.Fatalf("Failed to create service configuration: %v", err)
	}

	rtManager, err := artifactory.New(serviceConfig)
	if err != nil {
		log.Fatalf("Failed to create Artifactory manager: %v", err)
	}

	// Define the repository settings
	repoKey := "my-new-repo"
	repoConfig := services.NewLocalRepositoryBaseParams()
	repoConfig.Key = repoKey
	repoConfig.PackageType = "maven" // or any other type like "docker", "npm", etc.

	// Create the repository
	err = rtManager.CreateLocalRepository(repoConfig)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	fmt.Printf("Repository '%s' created successfully.\n", repoKey)
}
