package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/eiannone/keyboard"
	"google.golang.org/api/option"
)

// Movement represent chariot's movement data at a certain point in time
type Movement struct {
	Engine int `json:"engine,omitempty"`
	Wheel  int `json:"wheel,omitempty"`
}

var (
	client   *db.Client
	ref      *db.Ref
	movement Movement      = Movement{}
	duration time.Duration = 333 * time.Millisecond
)

func init() {
	sa := option.WithCredentialsFile("./key.json")
	config := firebase.Config{DatabaseURL: "https://chariot-rc.firebaseio.com/"}
	app, err := firebase.NewApp(context.Background(), &config, sa)
	if err != nil {
		log.Fatalf("Error initializing App: %v\n", err)
	}

	client, err = app.Database(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Database: %v\n", err)
	}

	ref = client.NewRef("chariot_1")
	fmt.Println(ref.Path)
}
func main() {
	resetMovement()
	go updateDatabse()

	keysEvents, err := keyboard.GetKeys(0)
	if err != nil {
		log.Fatalf("Error getting keyboard events: %v", err)
	}

	fmt.Println("Press ESC to quit")

	for true {
		event := <-keysEvents

		if event.Err != nil {
			panic(event.Err)
		}

		// Handle ESC
		if event.Key == keyboard.KeyEsc {
			fmt.Println("Shutdown")
			break
		}

		// Handle Ctrl + C
		if event.Rune == 0 && event.Key == 3 {
			shutdown()
		}

		event.Rune = unicode.ToLower(event.Rune)
		log.Printf("%q\n", event.Rune)

		switch event.Rune {
		case 'w':
			movement.Engine = 1
		case 's':
			movement.Engine = -1
		case 'a':
			movement.Wheel = -1
		case 'd':
			movement.Wheel = 1
		default:
			movement.Engine = 0
			movement.Wheel = 0
		}
	}
}

// UpdateDatabase sends the movement data to database
func updateDatabse() {
	for range time.Tick(duration) {
		fmt.Printf("engine: %d, wheel: %d\n", movement.Engine, movement.Wheel)

		err := ref.Set(context.Background(), map[string]interface{}{
			"wheel":  movement.Wheel,
			"engine": movement.Engine,
		})
		if err != nil {
			log.Fatalln("error updating database")
		}
	}

}

// Shutdown calls resetMovement closes the program.
func shutdown() {
	resetMovement()
	fmt.Println("Shutdown")
	os.Exit(0)
}

// ResetMovement resets the movement (both engine and wheel) to 0
func resetMovement() {
	err := ref.Set(context.Background(), map[string]interface{}{
		"wheel":  0,
		"engine": 0,
	})
	if err != nil {
		log.Fatalln("error updating database")
	}

	fmt.Println("Movement data resetted successfully.")
}
