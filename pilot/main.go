package main

import (
	"os"
	"context"
	"firebase.google.com/go"
	db "firebase.google.com/go/db"
	"github.com/eiannone/keyboard"
	"google.golang.org/api/option"
	"log"
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
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Database(context.Background())
	if err != nil {
		log.Fatalf("error initializing database: %v\n", err)
	}

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}

	log.Println("press ESC to quit")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		log.Printf("you pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			log.Printf("shutdown\n")
			break
		}

		if event.Rune == 0 && event.Key == 3 {
			log.Printf("shutdown\n")
			os.Exit(0)
		}

		go update(client, event.Rune)
	}
}

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

	client.NewRef("chariot_1").Set(context.Background(), &movement)
}
