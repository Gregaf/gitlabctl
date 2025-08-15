package commands

import (
	"fmt"
	"gregaf/gitlabctl/internal/services"
	"log/slog"
	"os"
)

type CommandHandler struct {
	gitlabService *services.GitlabService
	logger        *slog.Logger
}

func NewCommandHandler(gitlabService *services.GitlabService, logger *slog.Logger) *CommandHandler {
	return &CommandHandler{gitlabService: gitlabService, logger: logger}
}

type commandFunc func(args []string) (string, error)

func (ch *CommandHandler) Handle(cmd string, args []string) {
	help := "Add the subcommand 'edit' or 'list'!"

	ch.logger.Debug("Command handling", "args", args)

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(0)
	}

	commands := map[string]commandFunc{
		"verify": ch.CommandVerify,
	}

	handler, ok := commands[args[0]]
	if !ok {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(0)
	}

	out, err := handler(args)
	handleOutput(out, err)
}

func handleOutput(out string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, out)
	os.Exit(0)
}
