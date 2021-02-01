package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	regionID string

	title string
	body  string
	topic string
	token string

	verbose bool
)

var (
	firebaseApp     *firebase.App
	firestoreClient *firestore.Client
	messagingClient *messaging.Client
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&title, "title", "", "message title")
	flag.StringVar(&body, "body", "", "message body")
	flag.StringVar(&topic, "topic", "", "topic to send message to")
	flag.StringVar(&token, "token", "", "token of individual device to send message to")
	flag.BoolVar(&verbose, "verbose", false, "print extensive logs")

	opt := option.WithCredentialsFile("./key.json")

	var err error
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("notify: failed to initialize firebase app")
	}

	firestoreClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("notify: failed to initialize firestore: %v\n", err)
	}

	messagingClient, err = app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("notify: failed to initialize messaging")
	}
}

func main() {
	flag.Parse()

	if regionID == "" {
		log.Fatalln("notify: regionID is empty")
	}

	data := make(map[string]string)
	data["value"] = "dupa"

	msg := messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:  data,
		Token: token,
		Topic: topic,
	}

	b, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		log.Fatalln("notify: error marshalling message to json")
	}
	fmt.Println("notify: message to be sent:")
	fmt.Println(string(b))

	panic("lol")

	response, err := messagingClient.Send(context.Background(), &msg)
	if err != nil {
		log.Fatalf("notify: failed to send message: %v\n", err)
	}

	fmt.Println("notify: message sent, response: ", response)
}
