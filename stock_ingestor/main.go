package main

import (
	"log"

	"github.com/rohanchavan1918/stock_ingestor/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
