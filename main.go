package main

import (
	"api/app"
	"fmt"
)

func main() {
	fmt.Println("The task maintin api is starting")
	serverReady := make(chan bool)
	app := app.App{
		ServerReady: serverReady,
	}

	app.StartApp()
}
