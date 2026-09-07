package main

import (
	"fmt"
	"log"

	"github.com/Ha0cH/blogator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.ReadConfigJson()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	cfg.SetUser("Hao")

	cfg, err = config.ReadConfigJson()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}

	fmt.Printf("Config file: %+v\n", cfg)
}
