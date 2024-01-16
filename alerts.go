package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

var ApiKey = apikey
var BaseURL string

type AlertData struct {
	AlertType string `json:"alert_type"`
	Severity  string `json:"severity"`
	Agent     int    `json:"agent"`
	Message   string `json:"message"`
	Hostname  string `json:"hostname"`
	AgentID   string `json:"agent_id"`
	Client    string `json:"client"`
	Site      string `json:"site"`
}

func GetSeverity() string {
	switch Urgency {
	case "Low":
		return "info"
	case "Medium":
		return "warning"
	case "High":
		return "error"
	default:
		return "error" // This is a fallback in case Urgency has an unexpected value
	}
}

func getApiKeyFromAPI(token string, agentID string) (string, error) {
    apiUrl := "https://api.isfmb.com/api/v3/" + agentID + "/support/"

    client := &http.Client{}
    req, err := http.NewRequest("GET", apiUrl, nil)
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", "Token "+token)

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        // Read the response body to extract the API key
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return "", err
        }

        // Parse the JSON response to get the API key
        var responseData map[string]interface{}
        if err := json.Unmarshal(body, &responseData); err != nil {
            return "", err
        }

        apiKey, ok := responseData["support_token"].(string)
        if !ok {
            return "", fmt.Errorf("API key not found in response")
        }

        return apiKey, nil
    } else {
        return "", fmt.Errorf("Failed to retrieve API key: %d %s", resp.StatusCode, resp.Status)
    }
}

// Modify the SendAlert function to retrieve the API key before sending the POST request
func SendAlert() {
    // Fetch registry values
    reg, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\TacticalRMM`, registry.QUERY_VALUE)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer reg.Close()

    agentID, _, err := reg.GetStringValue("agentID")
    if err != nil {
        fmt.Println(err)
        return
    }

    token, _, err := reg.GetStringValue("Token") // Retrieve the "Token" value
    if err != nil {
        fmt.Println(err)
        return
    }

    hostname, err := os.Hostname()
    if err != nil {
        fmt.Println(err)
        return
    }

    baseURL, _, err := reg.GetStringValue("BaseURL")
    if err != nil {
        fmt.Println(err)
        return
    }

    BaseURL = baseURL

    // Retrieve the API key using the "Token"
    apiKey, err := getApiKeyFromAPI(token, agentID)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Prepare data
    data := AlertData{
        AlertType: "support",
        Severity:  GetSeverity(),
        Agent:     agentPK,
        Message:   fmt.Sprintf("%s has requested assistance. Their issue: "+Problem, hostname),
        Hostname:  hostname,
        AgentID:   agentID,
        Client:    "Demo",
        Site:      "Demo",
    }

    url := BaseURL + "/alerts/"

    // Send the POST request with the "X-API-KEY" header using the retrieved API key
    client := &http.Client{}
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
    if err != nil {
        fmt.Println(err)
        return
    }

    req.Header.Set("X-API-KEY", apiKey) // Set the "X-API-KEY" header using the retrieved API key
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    // Check the response
    if resp.StatusCode == http.StatusOK {
        fmt.Println("Alert created successfully:", resp)
    } else {
        fmt.Printf("Failed to create alert: %d %s\n", resp.StatusCode, resp.Status)
    }
}
