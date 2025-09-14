package rpc

import (
	"context"
	"fmt"
	"github.com/Brant-Liang/wallet-sign/chaindispatcher"
	"github.com/Brant-Liang/wallet-sign/config"
	"github.com/Brant-Liang/wallet-sign/leveldb"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"

	"github.com/Brant-Liang/wallet-sign/gen/go"
	"github.com/Brant-Liang/wallet-sign/hsm"
)

const MaxReceivedMessageSize = 1024 * 1024 * 64

type RpcServerConfig struct {
	GrpcHostname string
	GrpcPort     int
	KeyPath      string
	KeyName      string
	HsmEnable    bool
}

type RpcServer struct {
	conf      *config.Config
	HsmClient *hsm.HsmClient
	wallet.UnimplementedWalletServiceServer
	stopped    atomic.Bool
	gs         *grpc.Server
	ln         net.Listener
	dispatcher *chaindispatcher.ChainDispatcher
}

func NewRpcServer(cfg *config.Config) (*RpcServer, error) {
	rpcServer := &RpcServer{
		conf: cfg,
	}
	if cfg.HsmEnable {
		if cfg.KeyPath == "" || cfg.KeyName == "" {
			return nil, fmt.Errorf("hsm enabled but keyPath/keyName is empty")
		}
		hsmClient, hsmErr := hsm.NewHSMClient(context.Background(), cfg.KeyPath, cfg.KeyName)
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
		addr := fmt.Sprintf("%s:%d", s.conf.RpcServer.Host, s.conf.RpcServer.Port)
		log.Info("start rpc services", "addr", addr)

		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("Could not start tcp listener. ")
			return
		}
		s.ln = ln

		db, err := leveldb.NewKeyStore(s.conf.LevelDbPath)
		if err != nil {
			log.Error("Failed to create leveldb keystore", "err", err)
			return
		}

		dispatcher, err := chaindispatcher.NewChainDispatcher(s.conf, db)
		if err != nil {
			log.Error("new chain dispatcher fail", "err", err)
			return
		}
		s.dispatcher = dispatcher

		s.gs = grpc.NewServer(
			grpc.MaxRecvMsgSize(MaxReceivedMessageSize),
			grpc.ChainUnaryInterceptor(dispatcher.Interceptor),
		)

		wallet.RegisterWalletServiceServer(s.gs, dispatcher)

		reflection.Register(s.gs)

		log.Info("Grpc info", "port", s.conf.RpcServer.Port, "address", s.ln.Addr())

		go func() {
			if err := s.gs.Serve(s.ln); err != nil {
				// Serve 只有在 Stop/GracefulStop 或 ln 关闭/报错时返回
				log.Error("grpc serve returned", "err", err)
			}
		}()

		go func() {
			<-ctx.Done()
			s.internalStop("ctx.Done()")
		}()
	}(s)
	return nil
}

func (s *RpcServer) Stopped() bool {
	return s.stopped.Load()
}

func (s *RpcServer) internalStop(reason string) {
	log.Info("stopping rpc server", "reason", reason)
	if s.gs != nil {
		s.gs.GracefulStop()
	}
	if s.ln != nil {
		_ = s.ln.Close()
	}
	if s.HsmClient != nil && s.HsmClient.Gclient != nil {
		_ = s.HsmClient.Gclient.Close()
	}
}

func (s *RpcServer) Stop(ctx context.Context) error {
	s.stopped.Store(true)
	s.internalStop("Stop()")
	return nil
}
