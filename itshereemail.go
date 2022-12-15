package main

import (
	"io/ioutil"
	"log"
	"net/smtp"

	"github.com/fsnotify/fsnotify"
	"github.com/jhillyerd/enmime"
)

const (
	dir = "/path/to/directory"
	from = "sender@example.com"
	to = "recipient@example.com"
	smtpServer = "smtp.example.com:587"
	smtpUser = "sender@example.com"
	smtpPass = "PASSWORD"
)

func main() {
	// Create a new fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add the directory you want to watch to the watcher
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to receive events from the watcher
	events := watcher.Events

	// Start a separate goroutine to watch for events
	go func() {
		for {
			// Wait for an event
			event := <-events

			// Check if the event is a "create" event
			if event.Op&fsnotify.Create == fsnotify.Create {
				// The file has been created, you can now send it in an email
				sendEmail(event.Name)
			}
		}
	}()

	// Wait indefinitely
	select {}
}

func sendEmail(filename string) {
	// Read the file's contents
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new MIME message
	message := enmime.NewMessage()

	// Set the sender and recipient
	message.SetHeader("From", from)
	message.SetHeader("To", to)

	// Set the subject
	message.SetHeader("Subject", "New file detected")

	// Add the file as an attachment
	message.Attach(fileBytes, filename)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpUser, smtpPass, "smtp.example.com")
	err = smtp.SendMail(smtpServer, auth, from, []string{to}, message.Bytes())
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}
}
//This code will create a new fsnotify watcher and add the directory you want to monitor to it. It will then start a separate goroutine to watch for events on the watcher.
