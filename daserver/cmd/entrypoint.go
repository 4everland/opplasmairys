package main

import (
	"fmt"
	"github.com/4everland/opplasmairys"
	"github.com/4everland/opplasmairys/daserver"
	"golang.org/x/exp/slog"

	"github.com/urfave/cli/v2"
)

func StartDAServer(cliCtx *cli.Context) error {
	if err := CheckRequired(cliCtx); err != nil {
		return err
	}

	cfg := ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}

	slog.Info("Initializing Plasma DA server...")

	var store daserver.KVStore

	if cfg.FileStoreEnabled() {
		slog.Info("Using file storage", "path", cfg.FileStoreDirPath)
		store = NewFileStore(cfg.FileStoreDirPath)
	} else if cfg.S3Enabled() {
		slog.Info("Using S3 storage", "bucket", cfg.S3Config().Bucket)
		s3, err := NewS3Store(cfg.S3Config())
		if err != nil {
			return fmt.Errorf("failed to create S3 store: %w", err)
		}
		store = s3
	}

	//initial permaweb store
	if cfg.IrysEnabled() {
		var err error
		store, err = opplasmairys.NewDAStore(cfg.IrysConfig(), store)
		if err != nil {
			return fmt.Errorf("failed to create permaweb store: %w", err)
		}
	}
	server := daserver.NewDAServer(cliCtx.String(ListenAddrFlagName), cliCtx.Int(PortFlagName), store)

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start the DA server")
	} else {
		slog.Info("Started DA Server")
	}

	defer func() {
		if err := server.Stop(); err != nil {
			slog.Error("failed to stop DA server", "err", err)
		}
	}()

	BlockOnInterrupts()

	return nil
}
