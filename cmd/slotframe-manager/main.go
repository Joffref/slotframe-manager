package main

import (
	"github.com/Joffref/slotframe-manager/cmd/slotframe-manager/app"
)

func main() {
	if err := app.Execute(); err != nil {
		panic(err)
	}
}
