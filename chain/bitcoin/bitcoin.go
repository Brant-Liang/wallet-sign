package bitcoin

import (
	"context"
	"encoding/json"
	"github.com/Brant-Liang/wallet-sign/chain"
	"github.com/Brant-Liang/wallet-sign/config"
	wallet "github.com/Brant-Liang/wallet-sign/gen/go"
	"github.com/Brant-Liang/wallet-sign/hsm"
	"github.com/Brant-Liang/wallet-sign/leveldb"
	"github.com/Brant-Liang/wallet-sign/ssm"
	"github.com/ethereum/go-ethereum/log"
)

const ChainName = "Bitcoin"

type ChainAdaptor struct {
	signer    ssm.Signer
	db        *leveldb.Keys
	hsmClient *hsm.HsmClient
}

func NewChainAdapter(conf *config.Config, db *leveldb.Keys, hsmClient *hsm.HsmClient) (chain.IChainAdaptor, error) {
	return &ChainAdaptor{
		db:        db,
		hsmClient: hsmClient,
		signer:    ssm.NewEcdsaSigner(),
	}, nil
}

func (c ChainAdaptor) GetChainSignMethod(ctx context.Context, req *wallet.GetChainSignMethodRequest) (*wallet.GetChainSignMethodResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetChainSchema(ctx context.Context, req *wallet.GetChainSchemaRequest) (*wallet.GetChainSchemaResponse, error) {
	var vins []*Vin
	vins = append(vins, &Vin{
		Hash:   "",
		Index:  0,
		Amount: 0,
	})
	var vouts []*Vout
	vouts = append(vouts, &Vout{
		Address: "",
		Index:   0,
		Amount:  0,
	})
	bs := BitcoinSchema{
		RequestId: "0",
		Fee:       "0",
		Vins:      vins,
		Vouts:     vouts,
	}
	b, err := json.Marshal(bs)
	if err != nil {
		log.Error("marshal fail", "err", err)
	}
	return &wallet.GetChainSchemaResponse{
		Code:    wallet.ReturnCode_SUCCESS,
		Message: "get bitcoin sign schema success",
		Schema:  string(b),
	}, nil
}

func (c ChainAdaptor) CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairsExportPublicKeyListRequest) (*wallet.CreateKeyPairsExportPublicKeyListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) CreateKeyPairsWithAddresses(ctx context.Context, req *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) SignTransactionMessage(ctx context.Context, req *wallet.GetSignTransactionMessageRequest) (*wallet.GetSignTransactionMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildAndSignTransaction(ctx context.Context, req *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildAndSignBatchTransaction(ctx context.Context, req *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}
