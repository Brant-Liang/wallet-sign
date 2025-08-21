package rpc

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"

	"github.com/Brant-Liang/wallet-sign/gen/go"
	"github.com/Brant-Liang/wallet-sign/hsm"
	"github.com/Brant-Liang/wallet-sign/leveldb"
)

const MaxReceivedMessageSize = 1024 * 1024 * 4

type RpcServerConfig struct {
	GrpcHostname string
	GrpcPort     int
	KeyPath      string
	KeyName      string
	HsmEnable    bool
}

type RpcServer struct {
	*RpcServerConfig
	db        *leveldb.Keys
	HsmClient *hsm.HsmClient
	wallet.UnimplementedWalletServiceServer
	stopped atomic.Bool
}

func (s *RpcServer) Stop(ctx context.Context) error {
	s.stopped.Store(true)
	return nil
}

func (s *RpcServer) Stopped() bool {
	return s.stopped.Load()
}

func NewRpcServer(db *leveldb.Keys, config *RpcServerConfig) (*RpcServer, error) {
	var hsmClient *hsm.HsmClient
	var hsmErr error
	rpcServer := &RpcServer{
		RpcServerConfig: config,
		db:              db,
	}
	if config.HsmEnable {
		hsmClient, hsmErr = hsm.NewHSMClient(context.Background(), config.KeyPath, config.KeyName)
		if hsmErr != nil {
			log.Error("new hsm client fail", "err", hsmErr)
			return nil, hsmErr
		}
		rpcServer.HsmClient = hsmClient
	}
	return rpcServer, nil
}

func (s *RpcServer) Start(ctx context.Context) error {
	go func(s *RpcServer) {
		addr := fmt.Sprintf("%s:%d", s.GrpcHostname, s.GrpcPort)
		log.Info("start rpc services", "addr", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("Could not start tcp listener. ")
		}

		opt := grpc.MaxRecvMsgSize(MaxReceivedMessageSize)

		gs := grpc.NewServer(
			opt,
			grpc.ChainUnaryInterceptor(
				nil,
			),
		)
		reflection.Register(gs)

		wallet.RegisterWalletServiceServer(gs, s)

		log.Info("Grpc info", "port", s.GrpcPort, "address", listener.Addr())
		if err := gs.Serve(listener); err != nil {
			log.Error("Could not GRPC services")
		}
	}(s)
	return nil
}
