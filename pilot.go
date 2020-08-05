package main

import (
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

type Data struct {
        Engine int `json:"engine,omitempty"`
}

func main() {
	sa := option.WithCredentialsFile("./key.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("Error initializing Firebase app")
	}

	client, err := app.Database(context.Background())
	if err != nil {
		log.Fatalf("Error initializing RTDB")
	}

	ref := client.NewRef("chariot-rc").Child("chariot_1")

	var data Data
	err := ref.Get(context.Background(), &data)
}
