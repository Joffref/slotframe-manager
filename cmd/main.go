package main

import "github.com/Joffref/slotframe-manager/cmd/app"

func main() {
	if err := app.Execute(); err != nil {
		panic(err)
	}
}
