package main

import (
	"log/slog"
	"os"

	"github.com/amanbolat/zederr/internal/codegen/command"
)

func main() {
	slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       nil,
		ReplaceAttr: nil,
	}))

	err := command.NewRoot().Execute()
	if err != nil {
		os.Exit(1)
	}
}
