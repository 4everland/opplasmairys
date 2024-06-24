package main

import (
	"context"
	"golang.org/x/exp/slog"
	"os"

	"github.com/urfave/cli/v2"
)

var Version = "v0.0.1"

func main() {

	app := cli.NewApp()
	app.Flags = ProtectFlags(Flags)
	app.Version = FormatVersion(Version, "", "", "")
	app.Name = "da-server"
	app.Usage = "Plasma DA Storage Service"
	app.Description = "Service for storing plasma DA inputs"
	app.Action = StartDAServer

	ctx := WithInterruptBlocker(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		slog.Error("Application failed", "message", err)
	}

}
