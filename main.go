package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"systray"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var Urgency string
var Problem string

type SupportTicket struct {
	Name     string
	Email    string
	Problem  string
	Urgency  string
	Complete chan struct{}
}

var supportChannel = make(chan *SupportTicket)
var notifyIcon *walk.NotifyIcon
var supportWindowOpen bool = false

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	log.Println("onReady executed")

	systray.SetIcon(iconData)
	systray.SetTitle("Support Tray")
	systray.SetTooltip("iScreenFix RMM Support")

	menuSupport := systray.AddMenuItem("Support", "Request support")
	menuSupport.SetTooltip("Support Ticket")

	go func() {
		for {
			select {
			case <-menuSupport.ClickedCh:
				log.Println("Support clicked")
				OpenSupport()
			}
		}
	}()
}

func onExit() {
	if notifyIcon != nil {
		notifyIcon.Dispose()
	}
	log.Println("Exiting...")
}



func OpenSupport() {
    // Check if a window is already open
    if supportWindowOpen {
        log.Println("Support window already open.")
        return
    }
    supportWindowOpen = true

    var mw *walk.MainWindow

    icon, err := walk.NewIconFromResource("MYICON")
    if err != nil {
        log.Println("Error creating icon from resource:", err)
        return
    }

    var nameLineEdit, emailLineEdit, problemTextEdit *walk.TextEdit
    var urgencyComboBox *walk.ComboBox

    if err := (MainWindow{
        AssignTo: &mw,
        Title:    "Support Ticket",
        Icon:     icon,
        Size:     Size{400, 300},
        MinSize:  Size{400, 300},
        MaxSize:  Size{400, 300},
        Layout:   VBox{},
        Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
        Children: []Widget{
            Label{
                Text: "iScreenFix RMM Support",
                Font: Font{
                    Bold:      true,
                    Family:    "Segoe UI",
                    PointSize: 10,
                },
                TextAlignment: AlignCenter,
                MaxSize: Size{0, 20},
            },
            Label{Text: "Name:"},
            TextEdit{AssignTo: &nameLineEdit},
            Label{Text: "Email:"},
            TextEdit{AssignTo: &emailLineEdit},
            Label{Text: "Problem Description:"},
            TextEdit{AssignTo: &problemTextEdit, MinSize: Size{100, 50}},
            Label{Text: "Urgency:"},
            ComboBox{
                AssignTo: &urgencyComboBox,
                Model:    []string{"Low", "Medium", "High"},
            },
            PushButton{
                Text: "Submit",
                OnClicked: func() {
					name := nameLineEdit.Text()
					email := emailLineEdit.Text()
					problem := problemTextEdit.Text()
					urgency := urgencyComboBox.Text()
					Urgency = urgency
					Problem = problem

					hostname, err := os.Hostname()
					if err != nil {
						log.Println("Unable to retrieve hostname:", err)
						hostname = "UnknownDevice"
					}

					SendAlert()

					currentUser, err := user.Current()
					loggedInUser := "UnknownUser"
					if err == nil {
						loggedInUser = currentUser.Username
					}

					supportWindowOpen = false

					ticket := fmt.Sprintf("Received ticket from %s (%s): %s - Urgency: %s", name, email, problem, urgency)
					log.Println(ticket)

					subject := hostname + " has requested assistance"
					content := fmt.Sprintf("<b>Name:</b> %s<br><b>Email:</b> %s<br><b>Problem:</b> %s<br><b>Urgency:</b> %s<br><b>Logged in user:</b> %s<br><b>Device:</b> %s", name, email, problem, urgency, loggedInUser, hostname)
					TriggerEmail(subject, content)


					    log.Println("Attempting to show notification...")

if notifyIcon == nil {
    log.Println("Initializing new NotifyIcon...")
    notifyIcon, err = walk.NewNotifyIcon(mw)
    if err != nil {
        log.Println("Failed to create NotifyIcon:", err)
        return
    }
}


if err := notifyIcon.SetIcon(icon); err != nil {
    log.Println("Failed to set icon for NotifyIcon:", err)
    return
}

if err := notifyIcon.SetToolTip("Ticket Submitted!"); err != nil {
    log.Println("Failed to set tooltip for NotifyIcon:", err)
    return
}

if err := notifyIcon.SetVisible(true); err != nil {
    log.Println("Failed to set NotifyIcon visible:", err)
    return
}

if err := notifyIcon.ShowCustom("Support Ticket", "Your support ticket has been submitted!", icon); err != nil {
	log.Println("Failed to show info balloon for NotifyIcon:", err)
	return
}

log.Println("Notification should be visible now!")

// If you want the NotifyIcon itself to disappear after the balloon has been shown for a while (e.g., 10 seconds), 
// you can use the following, but this is optional:

time.AfterFunc(10*time.Second, func() {
    notifyIcon.SetVisible(false)
})


// Close the main window after some time.
time.AfterFunc(1*time.Second, func() {
    mw.Close()
})
				 },
            },
            Composite{
                Layout: HBox{},
                Children: []Widget{
                    HSpacer{},
                    PushButton{
                        Text: "Close",
                        OnClicked: func() {
                            mw.Close()
                        },
                    },
                },
            },
        },
    }).Create(); err != nil {
        log.Println("Error creating main window:", err)
        return
    }

    // Attach the handler to the Closing() event of the main window here
    mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
        supportWindowOpen = false
    })

    mw.Run()
}



