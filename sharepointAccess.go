package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const (
    tenantID     = "YOUR_TENANT_ID"
    clientID     = "YOUR_CLIENT_ID"
    clientSecret = "YOUR_CLIENT_SECRET"
    siteID       = "YOUR_SITE_ID"
    listID       = "YOUR_LIST_ID"
)

func getAccessToken() (string, error) {
    form := url.Values{}
    form.Add("grant_type", "client_credentials")
    form.Add("client_id", clientID)
    form.Add("client_secret", clientSecret)
    form.Add("resource", "https://graph.microsoft.com")

    resp, err := http.PostForm(fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenantID), form)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

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

func validateTokenClaims(accessToken string) (bool, error) {
    parts := strings.Split(accessToken, ".")
    if len(parts) != 3 {
        return false, fmt.Errorf("invalid token format")
    }

    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return false, err
    }

    var claims map[string]interface{}
    if err := json.Unmarshal(payload, &claims); err != nil {
        return false, err
    }

    if scp, ok := claims["scp"]; ok && scp != "" {
        return true, nil
    }
    if roles, ok := claims["roles"]; ok && len(roles.([]interface{})) > 0 {
        return true, nil
    }

    return false, fmt.Errorf("neither scp nor roles claim is present in the token")
}

func getSharePointList(accessToken string) error {
    client := &http.Client{}
    req, err := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/sites/%s/lists/%s/items", siteID, listID), nil)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer "+accessToken)
    req.Header.Set("Accept", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return err
    }

    // Process the list items
    fmt.Println(result)

    return nil
}

func main() {
    accessToken, err := getAccessToken()
    if err != nil {
        fmt.Println("Error getting access token:", err)
        return
    }

    isValid, err := validateTokenClaims(accessToken)
    if err != nil {
        fmt.Println("Error validating token claims:", err)
        return
    }

    if !isValid {
        fmt.Println("Token does not contain required claims")
        return
    }

    if err := getSharePointList(accessToken); err != nil {
        fmt.Println("Error getting SharePoint list:", err)
    }
}
