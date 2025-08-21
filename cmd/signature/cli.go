package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Brant-Liang/wallet-sign/common/cliapp"
	"github.com/Brant-Liang/wallet-sign/config"
	flags2 "github.com/Brant-Liang/wallet-sign/flags"
	"github.com/Brant-Liang/wallet-sign/leveldb"
	"github.com/Brant-Liang/wallet-sign/services/rpc"
	"github.com/ethereum/go-ethereum/log"
)

func runRpc(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	fmt.Println("running grpc services...")
	cfg := config.NewConfig(ctx)
	grpcServerCfg := &rpc.RpcServerConfig{
		GrpcHostname: cfg.RpcServer.Host,
		GrpcPort:     cfg.RpcServer.Port,
		KeyName:      cfg.KeyName,
		KeyPath:      cfg.CredentialsFile,
		HsmEnable:    cfg.HsmEnable,
	}
	db, err := leveldb.NewKeyStore(cfg.LevelDbPath)
	if err != nil {
		log.Error("new key store level db", "err", err)
	}
	return rpc.NewRpcServer(db, grpcServerCfg)
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
