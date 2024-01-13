package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "io/ioutil"
)

// TwilioSettings holds the configuration for Twilio API
type TwilioSettings struct {
    TwilioAccountSID   string   `json:"twilio_account_sid"`
    TwilioAuthToken    string   `json:"twilio_auth_token"`
    TwilioNumber       string   `json:"twilio_number"`
    SMSAlertRecipients []string `json:"sms_alert_recipients"`
}

// SMSData represents the data needed to send an SMS
type SMSData struct {
    To   string `json:"to"`
    Body string `json:"body"`
}

// FetchTwilioSettings fetches Twilio settings from the server
func FetchTwilioSettings() (TwilioSettings, error) {
    var settings TwilioSettings
    url := BaseURL + "/core/settings/"

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return settings, fmt.Errorf("error creating request: %v", err)
    }
    req.Header.Set("X-API-KEY", ApiKey)

    resp, err := client.Do(req)
    if err != nil {
        return settings, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
        return settings, fmt.Errorf("error decoding response: %v", err)
    }

    return settings, nil
}

// SendSMS sends an SMS using Twilio's API
func SendSMS(settings TwilioSettings, data SMSData) error {
    urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", settings.TwilioAccountSID)

    msgData := url.Values{}
    msgData.Set("To", data.To)
    msgData.Set("From", settings.TwilioNumber)
    msgData.Set("Body", data.Body)
    msgDataReader := *strings.NewReader(msgData.Encode())

    client := &http.Client{}
    req, err := http.NewRequest("POST", urlStr, &msgDataReader)
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }
    req.SetBasicAuth(settings.TwilioAccountSID, settings.TwilioAuthToken)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    // Read the response body for detailed error info
    responseBody, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        fmt.Println("Message sent successfully to:", data.To)
    } else {
        fmt.Printf("Failed to send message to %s. Status code: %d. Response: %s\n", data.To, resp.StatusCode, string(responseBody))
        return fmt.Errorf("failed to send message to %s: status code %d", data.To, resp.StatusCode)
    }

    return nil
}

// TriggerSMS formats and sends an SMS with the given details
// TriggerSMS function
func TriggerSMS(name, email, problem, urgency, loggedInUser, device string) {
    // Fetch Twilio settings
    twilioSettings, err := FetchTwilioSettings()
    if err != nil {
        fmt.Println("Error fetching Twilio settings:", err)
        return
    }

    // Check if Twilio settings are available and SMS recipients are defined
    if twilioSettings.TwilioAccountSID == "" || twilioSettings.TwilioAuthToken == "" || twilioSettings.TwilioNumber == "" || len(twilioSettings.SMSAlertRecipients) == 0 {
        fmt.Println("SMS configuration not provided, skipping SMS alert.")
        return
    }

    // Formatting the message as plain text
    smsContent := fmt.Sprintf("Sent from your Twilio trial account - %s has requested assistance:\nName: %s\nEmail: %s\nProblem: %s\nUrgency: %s\nLogged in user: %s\nDevice: %s",
        device, name, email, problem, urgency, loggedInUser, device)

    // Send SMS to each recipient
    for _, recipient := range twilioSettings.SMSAlertRecipients {
        smsData := SMSData{
            To:   recipient,
            Body: smsContent,
        }

        err = SendSMS(twilioSettings, smsData)
        if err != nil {
            fmt.Println("Error sending SMS to", recipient, ":", err)
        }
    }
}




