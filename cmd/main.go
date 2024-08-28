package main

import (
	"log"
	"test-task/initializer"
)

func main() {
	app, err := initializer.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	err = app.Engine.Run(":" + app.Config.Port)
	if err != nil {
		log.Fatal("error running server: ", err)
	}
}
