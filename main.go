package main

import (
	"fmt"
	"log"
	"mooshi-1/aggregator/internal/config"
	"os"
)

func main() {
	cfgG := config.ReadConfig()
	cfgG.SetUser("Mooshi-1")

	contents, err := os.ReadFile(cfgG.Path)
	if err != nil {
		log.Fatal("error reading config after writing")
	}

	fmt.Print(string(contents))

}
