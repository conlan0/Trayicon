# Trayicon
In alerts.go you will need to replace your_api_key with a valid api key to your rmm. After that everything will compile and work as expected.

What the program does:
1. starts a system tray icon that has a submenu with a support item
2. that support item opens a support ticket form window with the fields (Name, email, problem description and a combobox with 3 choices for urgency (Low, Medium and High))
3. Once the form is filled out and submitted it triggers a toast balloon notifications letting the end user know that their ticket it submitted, it checks the /core/settings endpoint for email_alert_recipients and uses the user at 0, it also uses /core/settings for the smtp_from_email, smtp_host, smtp_host_password, smtp_host_user and smtp_port for sending the email ticket
4. When sending the emails and alerts the application also uses the registry to determine the hostname, and then in Software/TacticalRmm it uses agentid, agentpk and baseurl
5. It then uses the /alerts endpoint to post an alert to the dashboard with the type "support" and the urgency aligns with the severity (High:error, Medium:warning, Low:info) and then closes the support ticket window and continues resting in the system tray until used again
