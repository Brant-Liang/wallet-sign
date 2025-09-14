package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Brant-Liang/wallet-sign/common/cliapp"
	"github.com/Brant-Liang/wallet-sign/config"
	flags2 "github.com/Brant-Liang/wallet-sign/flags"
	"github.com/Brant-Liang/wallet-sign/services/rpc"
	"github.com/ethereum/go-ethereum/log"
)

func runRpc(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	fmt.Println("running grpc services...")
	cfgFile := ctx.String("config")
	if cfgFile == "" {
		cfgFile = "config.yml"
	}
	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		log.Error("config.NewConfig error", "error", err)
		return nil, err
	}

	srv, err := rpc.NewRpcServer(cfg)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func NewCli(GitCommit string, GitData string) *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              "v0.0.1-beta",
		Description:          "wallet sign rpc service",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "rpc",
				Flags:       flags,
				Description: "Run rpc services",
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
