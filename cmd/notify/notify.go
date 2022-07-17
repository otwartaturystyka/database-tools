// Package notify implements sending push notifications to the mobile app.
package notify

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/opentouristics/database-tools/readers"
	"google.golang.org/api/option"
)

var (
	title string
	body  string
	topic string
	token string
)

var messagingClient *messaging.Client

func init() {
	log.SetFlags(0)
}

func InitFirebase() error {
	flag.StringVar(&title, "title", "", "message title")
	flag.StringVar(&body, "body", "", "message body")
	flag.StringVar(&topic, "topic", "", "topic to send message to")
	flag.StringVar(&token, "token", "", "token of individual device to send message to")

	opt := option.WithCredentialsFile("./key.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("initialize firebase app: %v", err)
	}

	messagingClient, err = app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("notify: failed to initialize messaging: %v", err)
	}

	return nil
}

// Notify sends a push notification to users of the app who have regionID set as
// their default region.
func Notify(regionID string, verbose bool) error {
	if regionID == "" {
		return fmt.Errorf("regionID is empty")
	}

	data := make(map[string]string)
	data["title"] = title
	data["body"] = body

	msg := messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Token: token,
		Topic: topic,
	}

	msgJSON, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal message to JSON: %v", err)
	}
	fmt.Println("message to be sent:", msgJSON)

	confirmed, err := readers.AskForConfirmation(os.Stdin, os.Stdout, "send the message?", false)
	if err != nil {
		return fmt.Errorf("ask for confirmation: %v", err)
	}

	if !confirmed {
		return errors.New("canceled")
	}

	messagingResponse, err := messagingClient.Send(context.Background(), &msg)
	if err != nil {
		return fmt.Errorf("send message: %v", err)
	}

	fmt.Println("message sent, messagingResponse: ", messagingResponse)

	return nil
}
