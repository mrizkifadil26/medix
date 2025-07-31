package main

import (
	"log"

	"github.com/mrizkifadil26/medix/classifier"
	"github.com/mrizkifadil26/medix/utils/cli"
	"github.com/mrizkifadil26/medix/utils/config"
)

type CLI struct {
	ConfigPath string `flag:"config" help:"Path to classifier config file"`
}

func main() {
	var args CLI
	cli.Parse(&args)

	var cfg classifier.Config
	if err := config.Parse(args.ConfigPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := classifier.Run(cfg); err != nil {
		log.Fatalf("classifier error: %v", err)
	}
}
