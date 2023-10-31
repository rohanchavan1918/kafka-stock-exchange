package main

import (
	"log"

	"github.com/rohanchavan1918/order_processor/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
