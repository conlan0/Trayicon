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

	SendEmail(settings, subject, body)
}

func SendEmail(settings Settings, subject string, body string) {
	to := settings.EmailRecipients
	from := settings.SMTPFromEmail
	// subject := "Test Email"
	// body := "This is a test email sent from the RMM system."

	header := make(map[string]string)
	header["From"] = from
	header["To"] = to[0] // Assuming at least one recipient, use the first recipient here
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	serverName := fmt.Sprintf("%s:%d", settings.SMTPHost, settings.SMTPPort)
	auth := smtp.PlainAuth("", settings.SMTPHostUser, settings.SMTPHostPassword, settings.SMTPHost)

	err := smtp.SendMail(serverName, auth, from, to, []byte(message))
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
