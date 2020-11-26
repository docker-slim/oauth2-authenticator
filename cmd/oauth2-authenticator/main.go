package main

import (
	"fmt"

	"github.com/docker-slim/oauth2-authenticator/pkg/authenticator"
)

func main() {
	app := authenticator.New(nil)
	app.Run()

	if app.Input.Verbose {
		fmt.Println("ACCESS TOKEN:")
	}

	fmt.Printf("%s\n", app.AccessToken)

	if app.Input.ShowUser {
		if app.Input.Verbose {
			fmt.Println("USER INFO:")
		}

		fmt.Printf("%s\n", app.UserJSON())
	}
}

