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

var ApiKey = "your_api_key"
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

	agentPKStr, _, err := reg.GetStringValue("AgentPK")
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

	// Convert agentPKStr to int
	agentPK, err := strconv.Atoi(agentPKStr)
	if err != nil {
		fmt.Println("Failed to convert AgentPK to int:", err)
		return
	}

	// Prepare data
	data := AlertData{
		AlertType: "custom",
		Severity:  GetSeverity(),
		Agent:     agentPK,
		Message:   fmt.Sprintf("%s has requested assistance. Their issue: "+Problem, hostname),
		Hostname:  hostname,
		AgentID:   agentID,
		Client:    "Demo",
		Site:      "Demo",
	}

	url := BaseURL + "/alerts/"

	// Send the POST request
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("X-API-KEY", ApiKey)
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
