package main

import (
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"syscall"
	"systray"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	systray.Run(onReady, onExit)
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "TacticalRMMService",
		DisplayName: "Tactical RMM Service",
		Description: "Tactical Remote Monitoring and Management system",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func onReady() {
	iconData, err := ioutil.ReadFile("favicon.ico")
	if err != nil {
		log.Fatal(err)
	}
	systray.SetIcon(iconData)
	systray.SetTitle("Tactical RMM")
	systray.SetTooltip("Tactical RMM")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	virusscan := systray.AddMenuItem("Virus Scan", "Scans Computer for malware")
	go func() {
		for {
			<-virusscan.ClickedCh
			cmd := exec.Command("C:\\Program Files (x86)\\Kaspersky Lab\\KES.11.10.0\\avp.exe", "SCAN", "/all")
			virusscan.SetTooltip("This may slow down your computer temporarily.")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			err := cmd.Start()
			if err != nil {
				log.Fatalf("Error starting program: %v", err)
			}
			err = cmd.Wait()
			if err != nil {
				log.Fatalf("Error waiting for program to finish: %v", err)
			}
			log.Println("Virus scan completed successfully")
		}
	}()

	resourceMonitor := systray.AddMenuItem("Resource Monitor", "View system resource usage")
	go func() {
		for {
			<-resourceMonitor.ClickedCh
			cmd := exec.Command("resmon.exe")
			err := cmd.Start()
			if err != nil {
				log.Fatalf("Error starting program: %v", err)
			}
			log.Println("Resource Monitor launched successfully")
		}
	}()

	// Add a new menu item to send email
	mEmail := systray.AddMenuItem("Support", "Request support")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		log.Println("Quit now...")
	}()

	// Handle the "Send Email" menu item click
	go func() {
		for {
			<-mEmail.ClickedCh
			desktopName, _ := os.Hostname()
			sendEmail("from@example.com", "to@example.com", "smtp.example.com", "username", "password", desktopName)
		}
	}()
}

func onExit() {
	// clean up here
	log.Println("Exiting...")
}

func sendEmail(from string, to string, smtpServer string, smtpUsername string, smtpPassword string, desktopName string) {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + desktopName + " has requested assistance\r\n" +
		"\r\n" +
		"This Desktop needs assistance.\r\n")

	err := smtp.SendMail(smtpServer+":587", auth, from, []string{to}, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Request sent.")
}
