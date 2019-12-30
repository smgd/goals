package main

import "goals/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
