package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-toast/toast"
)

func makeToast() {
	tmpDir := os.TempDir()
	iconFilePath := filepath.Join(tmpDir, "temp_icon.ico")
	err := ioutil.WriteFile(iconFilePath, iconData, 0644)
	if err != nil {
		log.Fatalf("Failed to write icon data to temp file: %v", err)
	}
	defer os.Remove(iconFilePath)

	notification := toast.Notification{
		AppID:   "Tatcical RMM",
		Title:   "RMM Support Ticket",
		Message: "Your Support Ticket has been submitted!",
		Icon:    iconFilePath,
		Actions: []toast.Action{
			{"protocol", "Dismiss", ""},
		},
	}

	err = notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
