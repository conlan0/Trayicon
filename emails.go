package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/smtp"
)

type Settings struct {
    SMTPHost         string   `json:"smtp_host"`
    SMTPPort         int      `json:"smtp_port"`
    SMTPHostUser     string   `json:"smtp_host_user"`
    SMTPHostPassword string   `json:"smtp_host_password"`
    SMTPFromEmail    string   `json:"smtp_from_email"`
    EmailRecipients  []string `json:"email_alert_recipients"`
}

func TriggerEmail(subject string, body string) {
    url := BaseURL + "/core/settings/"

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    req.Header.Set("X-API-KEY", ApiKey)

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var settings Settings
    if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
        panic(err)
    }

    // Check if all SMTP settings are provided
    if settings.SMTPHost == "" || settings.SMTPPort == 0 || 
       settings.SMTPHostUser == "" || settings.SMTPHostPassword == "" || 
       settings.SMTPFromEmail == "" || len(settings.EmailRecipients) == 0 {
        fmt.Println("SMTP configuration not complete, skipping email alert.")
        return
    }

    SendEmail(settings, subject, body)
}

func SendEmail(settings Settings, subject string, body string) {
    from := settings.SMTPFromEmail

    header := make(map[string]string)
    header["From"] = from
    header["Subject"] = subject
    header["MIME-Version"] = "1.0"
    header["Content-Type"] = "text/html; charset=\"utf-8\""

    serverName := fmt.Sprintf("%s:%d", settings.SMTPHost, settings.SMTPPort)
    auth := smtp.PlainAuth("", settings.SMTPHostUser, settings.SMTPHostPassword, settings.SMTPHost)

    for _, to := range settings.EmailRecipients {
        header["To"] = to
        message := ""
        for k, v := range header {
            message += fmt.Sprintf("%s: %s\r\n", k, v)
        }
        message += "\r\n" + body

        err := smtp.SendMail(serverName, auth, from, []string{to}, []byte(message))
        if err != nil {
            fmt.Printf("Failed to send email to %s: %v\n", to, err)
        } else {
            fmt.Printf("Email sent successfully to %s\n", to)
        }
    }
}

