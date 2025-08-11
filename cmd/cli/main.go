package main

import (
	"flag"
	"fmt"
	"gregaf/gitlabctl/internal/commands"
	"gregaf/gitlabctl/internal/config"
	"gregaf/gitlabctl/internal/services"
	"log/slog"
	"os"
)

var (
	_flagConfig  = flag.String("config", "config.json", "Path to the configuration file")
	_flagVerbose = flag.Bool("verbose", false, "Turn on verbose logging")
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please provide a subcommand (e.g., 'verify').")
		os.Exit(1)
	}

	config, err := config.LoadConfig(*_flagConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	gitlabService := services.NewGitlabService(config, logger)
	commandHandler := commands.NewCommandHandler(gitlabService, logger)

	commandHandler.Handle("cmd", args)
}
