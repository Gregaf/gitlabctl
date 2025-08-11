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

func (ch *CommandHandler) Handle(cmd string, args []string) {
	help := "Add the subcommand 'edit' or 'list'!"

	ch.logger.Debug("Command handling", "args", args)

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(0)
	}
	switch args[0] {
	case "verify":
		ch.CommandVerify(args)
	default:
		fmt.Fprintln(os.Stderr, help)
		os.Exit(0)
	}

}
