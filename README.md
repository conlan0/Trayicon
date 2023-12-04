# Trayicon
In alerts.go you will need to replace your_api_key with a valid api key to your rmm and depending on if you recompiled your frontend with the changes you might need to change the AlertType in alerts.go from "support" to "custom". After that everything will compile and work as expected.

What the program does:
1. starts a system tray icon that has a submenu with a support item
2. that support item opens a support ticket form window with the fields (Name, email, problem description and a combobox with 3 choices for urgency (Low, Medium and High))
3. Once the form is filled out and submitted it triggers a toast balloon notifications letting the end user know that their ticket is submitted, it checks the /core/settings/ endpoint for email_alert_recipients and uses the user at 0, it also uses /core/settings/ for the smtp_from_email, smtp_host, smtp_host_password, smtp_host_user and smtp_port for sending the email ticket
4. When sending the emails and alerts the application also uses the registry to determine the hostname, and then in Software/TacticalRmm it uses agentid, agentpk and baseurl
5. It then uses the /alerts/ endpoint to post an alert to the dashboard with the type "support" and the urgency aligns with the severity (High:error, Medium:warning, Low:info) and then closes the support ticket window and continues resting in the system tray until used again


# RMM Additions
![image](https://github.com/conlan0/Trayicon/assets/87742085/2580a003-baee-4107-a46b-7084cfc21ce7)

I have also made small additions to the api and frontend of the rmm. This just includes an Added alert type of support in api/tacticalrmm/tacticalrmm/constants.py and a new style and alert type check to apply that style in src/components/modals/alerts/AlertsOverview.vue
