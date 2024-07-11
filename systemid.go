//URL: https://<your-artifactory-url>/artifactory/api/storage/my-repo/path/to/artifact?properties=version=1.0.0;buildNumber=123
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

const (
    jfrogURL       = "https://<your-artifactory-url>/artifactory"
    apiKey         = "<your-api-key>"
    repoKey        = "my-repo"
    artifactPath   = "path/to/artifact"
    newRepoJSON    = `{
        "rclass": "local",
        "packageType": "generic",
        "description": "My new repository",
        "repoLayoutRef": "simple-default"
    }`
    properties     = "version=1.0.0;buildNumber=123"
)

func main() {
    // Step 1: Create Repository
    if err := createRepository(repoKey, newRepoJSON); err != nil {
        fmt.Println("Error creating repository:", err)
        return
    }
    fmt.Println("Repository created successfully.")

    // Step 2: Add Properties to Artifact
    if err := addPropertiesToArtifact(repoKey, artifactPath, properties); err != nil {
        fmt.Println("Error adding properties to artifact:", err)
        return
    }
    fmt.Println("Properties added to artifact successfully.")
}

func createRepository(repoKey, jsonStr string) error {
    url := fmt.Sprintf("%s/api/repositories/%s", jfrogURL, repoKey)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonStr)))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-JFrog-Art-Api", apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        body, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    return nil
}

func addPropertiesToArtifact(repoKey, artifactPath, properties string) error {
    url := fmt.Sprintf("%s/api/storage/%s/%s?properties=%s", jfrogURL, repoKey, artifactPath, properties)
    req, err := http.NewRequest("PUT", url, nil)
    if err != nil {
        return err
    }
    req.Header.Set("X-JFrog-Art-Api", apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        body, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    return nil
}
