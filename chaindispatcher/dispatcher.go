package chaindispatcher

import (
	"context"
	"fmt"
	"github.com/Brant-Liang/wallet-sign/chain"
	"github.com/Brant-Liang/wallet-sign/chain/bitcoin"
	"github.com/Brant-Liang/wallet-sign/chain/ethereum"
	"github.com/Brant-Liang/wallet-sign/chain/solana"
	"github.com/Brant-Liang/wallet-sign/config"
	wallet "github.com/Brant-Liang/wallet-sign/gen/go"
	"github.com/Brant-Liang/wallet-sign/hsm"
	"github.com/Brant-Liang/wallet-sign/leveldb"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
)

const (
	AccessToken string = "DappLinkTheWeb3202402290001"
	WalletKey   string = "DappLinkWalletServicesRiskKeyxxxxxxxKey"
	RisKKey     string = "DappLinkWalletServicesRiskKeyxxxxxxxKey"
)

type CommonRequest interface {
	GetConsumerToken() string
	GetChainName() string
}

type ChainName = string
type ChainDispatcher struct {
	registry map[ChainName]chain.IChainAdaptor
}

type CommonResponse = wallet.GetChainSignMethodResponse

func NewChainDispatcher(conf *config.Config, db *leveldb.Keys) (*ChainDispatcher, error) {
	dispatcher := ChainDispatcher{
		registry: make(map[ChainName]chain.IChainAdaptor),
	}
	chainAdaptorFactoryMap := map[string]func(conf *config.Config, db *leveldb.Keys, hsm *hsm.HsmClient) (chain.IChainAdaptor, error){
		bitcoin.ChainName:  bitcoin.NewChainAdapter,
		ethereum.ChainName: ethereum.NewChainAdapter,
		solana.ChainName:   solana.NewChainAdapter,
	}
	supportedChainNames := []string{
		bitcoin.ChainName,
		ethereum.ChainName,
		solana.ChainName,
	}

	var hsmClient *hsm.HsmClient
	if conf.HsmEnable {
		var err error
		hsmClient, err = hsm.NewHSMClient(context.Background(), conf.KeyPath, conf.KeyName)
		if err != nil {
			return nil, fmt.Errorf("new hsm client fail: %w", err)
		}
	}

	for _, chainName := range conf.Chains {
		if factory, ok := chainAdaptorFactoryMap[chainName]; ok {
			adaptor, err := factory(conf, db, hsmClient)
			if err != nil {
				log.Error("chain adapter factory err:", err, "chain:", chainName)
			}
			dispatcher.registry[chainName] = adaptor
		} else {
			log.Error("unsupported chain:", chainName, "supported:", supportedChainNames)
		}
	}
	return &dispatcher, nil
}

func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]
	consumerToken := req.(CommonRequest).GetConsumerToken()
	chainName := req.(CommonRequest).GetChainName()
	log.Info(method, "chain", chainName, "req", req, "consumerToken", consumerToken)
	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *ChainDispatcher) preHandler(req interface{}) *CommonResponse {
	consumerToken := req.(CommonRequest).GetConsumerToken()
	log.Debug("consumer token", "consumerToken", consumerToken, "req", req)
	if consumerToken != AccessToken {
		return &CommonResponse{
			Code:    wallet.ReturnCode_ERROR,
			Message: "consumer token is error",
		}
	}
	chainName := req.(CommonRequest).GetChainName()
	log.Debug("chain name", "chain", chainName, "req", req)
	if _, ok := d.registry[chainName]; !ok {
		return &CommonResponse{
			Code:    wallet.ReturnCode_ERROR,
			Message: "unsupported chain",
		}
	}
	return nil
}

func (d *ChainDispatcher) GetChainSignMethod(ctx context.Context, req *wallet.GetChainSignMethodRequest) (*wallet.GetChainSignMethodResponse, error) {
	resp := d.preHandler(req)
	if resp != nil {
		return &wallet.GetChainSignMethodResponse{
			Code:    wallet.ReturnCode_ERROR,
			Message: "unsupported method",
		}, nil
	}
	return d.registry[req.ChainName].GetChainSignMethod(ctx, req)
}

func (d *ChainDispatcher) GetChainSchema(ctx context.Context, req *wallet.GetChainSchemaRequest) (*wallet.GetChainSchemaResponse, error) {
	resp := d.preHandler(req)
	if resp != nil {
		return &wallet.GetChainSchemaResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[req.ChainName].GetChainSchema(ctx, req)
}

func (d *ChainDispatcher) CreateKeyPairsExportPublicKeyList(ctx context.Context, request *wallet.CreateKeyPairsExportPublicKeyListRequest) (*wallet.CreateKeyPairsExportPublicKeyListResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.CreateKeyPairsExportPublicKeyListResponse{
			Code: resp.Code,
			Msg:  resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].CreateKeyPairsExportPublicKeyList(ctx, request)
}

func (d *ChainDispatcher) CreateKeyPairsWithAddresses(ctx context.Context, request *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.CreateKeyPairsWithAddressesResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].CreateKeyPairsWithAddresses(ctx, request)
}

func (d *ChainDispatcher) SignTransactionMessage(ctx context.Context, request *wallet.GetSignTransactionMessageRequest) (*wallet.GetSignTransactionMessageResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.GetSignTransactionMessageResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].SignTransactionMessage(ctx, request)
}

func (d *ChainDispatcher) BuildAndSignTransaction(ctx context.Context, request *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.BuildAndSignTransactionResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	//txReqJsonByte, err := base64.StdEncoding.DecodeString(request.TxBase64Body)
	//if err != nil {
	//	return &wallet.BuildAndSignTransactionResponse{
	//		Code:    wallet.ReturnCode_ERROR,
	//		Message: "decode base64 string fail",
	//	}, nil
	//}
	//RiskKeyHash := crypto.Keccak256(append(txReqJsonByte, []byte(RisKKey)...))
	//RistKeyHashStr := hexutils.BytesToHex(RiskKeyHash)
	//if RistKeyHashStr != request.RiskKeyHash {
	//	return &wallet.BuildAndSignTransactionResponse{
	//		Code:    wallet.ReturnCode_ERROR,
	//		Message: "riskKey hash check Fail",
	//	}, nil
	//}
	//WalletKeyHash := crypto.Keccak256(append(txReqJsonByte, []byte(WalletKey)...))
	//WalletKeyHashStr := hexutils.BytesToHex(WalletKeyHash)
	//if WalletKeyHashStr != request.WalletKeyHash {
	//	return &wallet.BuildAndSignTransactionResponse{
	//		Code:    wallet.ReturnCode_ERROR,
	//		Message: "wallet key hash Check Fail",
	//	}, nil
	//}
	return d.registry[request.ChainName].BuildAndSignTransaction(ctx, request)
}

func (d *ChainDispatcher) BuildAndSignBatchTransaction(ctx context.Context, request *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &wallet.BuildAndSignBatchTransactionResponse{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	return d.registry[request.ChainName].BuildAndSignBatchTransaction(ctx, request)
}
