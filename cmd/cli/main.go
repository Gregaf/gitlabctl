package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"gregaf/gitlabctl/internal/commands"
	"gregaf/gitlabctl/internal/config"
	"gregaf/gitlabctl/internal/services"
	"gregaf/gitlabctl/internal/utils"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var (
	_flagConfig   = flag.String("config", "config.json", "Path to the configuration file")
	_flagVerbose  = flag.Bool("verbose", false, "Turn on verbose logging")
	_flagInsecure = flag.Bool("insecure", false, "Turn off TLS")
)

// TODO: Shift logic down into command handler, main shouldnt need much logic.
func main() {
	flag.Parse()
	args := flag.Args()

	outputStream := utils.GetOutputStream(*_flagVerbose, os.Stderr, io.Discard)

	logger := slog.New(slog.NewTextHandler(outputStream, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please provide a subcommand (e.g., 'verify').")
		os.Exit(1)
	}

	config, err := config.LoadConfig(*_flagConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: *_flagInsecure,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	gitlabService := services.NewGitlabService(client, config, logger)
	commandHandler := commands.NewCommandHandler(gitlabService, logger)

	commandHandler.Handle("cmd", args)
}
