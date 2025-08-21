package main

import (
	"context"
	"github.com/Brant-Liang/wallet-sign/common/opio"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

var (
	GitCommit = ""
	GitDate   = ""
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelInfo, true)))
	app := NewCli(GitCommit, GitDate)
	ctx := opio.WithInterruptBlocker(context.Background())
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
