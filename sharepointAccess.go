package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const (
    tenantID       = "YOUR_TENANT_ID"
    clientID       = "YOUR_CLIENT_ID"
    clientSecret   = "YOUR_CLIENT_SECRET"
    resource       = "https://graph.microsoft.com"
    scope          = "https://graph.microsoft.com/.default"
    authorityHost  = "https://login.microsoftonline.com"
    siteID         = "YOUR_SITE_ID"
)

func getAccessToken() (string, error) {
    // Construct the request URL
    tokenURL := fmt.Sprintf("%s/%s/oauth2/v2.0/token", authorityHost, tenantID)

    // Create the request payload
    data := url.Values{}
    data.Set("client_id", clientID)
    data.Set("client_secret", clientSecret)
    data.Set("grant_type", "client_credentials")
    data.Set("scope", scope)

    // Create a new HTTP request
    req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Parse the JSON response
    var tokenResp map[string]interface{}
    if err := json.Unmarshal(body, &tokenResp); err != nil {
        return "", err
    }

    // Extract the access token
    accessToken := tokenResp["access_token"].(string)
    return accessToken, nil
}

func validateTokenClaims(accessToken string) (map[string]interface{}, error) {
    parts := strings.Split(accessToken, ".")
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid token format")
    }

    // Decode the payload
    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return nil, err
    }

    // Unmarshal JSON
    var claims map[string]interface{}
    if err := json.Unmarshal(payload, &claims); err != nil {
        return nil, err
    }

    return claims, nil
}

func main() {
    // Get access token
    accessToken, err := getAccessToken()
    if err != nil {
        fmt.Println("Error getting access token:", err)
        return
    }

    // Validate token claims
    claims, err := validateTokenClaims(accessToken)
    if err != nil {
        fmt.Println("Error validating token claims:", err)
        return
    }

    // Print scp and roles claims
    scpClaim, scpExists := claims["scp"]
    rolesClaim, rolesExists := claims["roles"]

    if scpExists {
        fmt.Println("scp:", scpClaim)
    } else {
        fmt.Println("scp claim is not present")
    }

    if rolesExists {
        fmt.Println("roles:", rolesClaim)
    } else {
        fmt.Println("roles claim is not present")
    }
}
