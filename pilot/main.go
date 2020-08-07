package main

import (
	"context"
	"log"
	"os"
	"unicode"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/eiannone/keyboard"
	"google.golang.org/api/option"
)

type movement struct {
	Engine int `json:"engine,omitempty"`
	Wheel  int `json:"wheel,omitempty"`
}

func main() {
	sa := option.WithCredentialsFile("./key.json")
	config := firebase.Config{DatabaseURL: "https://chariot-rc.firebaseio.com/"}
	app, err := firebase.NewApp(context.Background(), &config, sa)
	if err != nil {
		log.Fatalf("Error initializing App: %v\n", err)
	}

	client, err := app.Database(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Database: %v\n", err)
	}

	keysEvents, err := keyboard.GetKeys(0)
	if err != nil {
		panic(err)
	}

	log.Println("Press ESC to quit")
	
	for true {
		event := <-keysEvents

		if event.Err != nil {
			panic(event.Err)
		}

		// Handle ESC
		if event.Key == keyboard.KeyEsc {
			log.Println("Shutdown")
			break
		}
		
		// Handle Ctrl + C
		if event.Rune == 0 && event.Key == 3 {
			log.Println("Shutdown")
			os.Exit(0)
		}

		event.Rune = unicode.ToLower(event.Rune)
		log.Printf("%q\n", event.Rune)


		go update(client, event.Rune)
	}
}

// Determines the move and sends it to
func update(client *db.Client, char rune) {
	movement := new(movement)
	
	switch char {
	case 'w':
		movement.Engine = 1
	case 's':
		movement.Engine = -1
	case 'a':
		movement.Wheel = -1
	case 'd':
		movement.Wheel = 1
	}
	
	
	ref := client.NewRef("chariot_1")
	err := ref.Set(context.Background(), &movement)
	if err != nil {
		log.Fatalf("error updating database")
	}
}
