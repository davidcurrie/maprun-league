package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/davidcurrie/maprun-league/internal/config"
	"github.com/davidcurrie/maprun-league/internal/league"
	"github.com/davidcurrie/maprun-league/internal/publisher"
)

type CLI struct {
	Config string `help:"Path to the configuration file" type:"path" default:"config.yaml"`
}

func main() {
	var cli CLI
	kong.Parse(&cli)

	cfg, err := config.Load(cli.Config)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Println("Processing league results...")
	html, err := league.ProcessLeague(cfg)
	if err != nil {
		log.Fatalf("Error processing league: %v", err)
	}

	fmt.Println(html)

	if cfg.DrupalPublisher.Enabled {
		log.Printf("Publishing results to %s", cfg.DrupalPublisher.URL)
		err := publisher.Publish(html, cfg.DrupalPublisher.URL, cfg.DrupalPublisher.Username, cfg.DrupalPublisher.Password)
		if err != nil {
			log.Fatalf("Error publishing results: %v", err)
		}
		log.Println("Successfully published results.")
	} else {
		log.Println("Publisher is disabled. Writing to results.html")
		if err := os.WriteFile("results.html", []byte(html), 0644); err != nil {
			log.Fatalf("Failed to write results.html: %v", err)
		}
	}
}
