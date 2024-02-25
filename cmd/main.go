package main

import (
	"log/slog"
	"os"

	"github.com/amanbolat/zederr/cmd/gen"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "zederr"
	app.Usage = "A tool to work with standardized errors."
	app.Commands = []*cli.Command{
		&gen.CmdGen,
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
