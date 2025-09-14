package chain

import (
	"context"
	wallet "github.com/Brant-Liang/wallet-sign/gen/go"
)

type IChainAdaptor interface {
	GetChainSignMethod(ctx context.Context, req *wallet.GetChainSignMethodRequest) (*wallet.GetChainSignMethodResponse, error)
	GetChainSchema(ctx context.Context, req *wallet.GetChainSchemaRequest) (*wallet.GetChainSchemaResponse, error)
	CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairsExportPublicKeyListRequest) (*wallet.CreateKeyPairsExportPublicKeyListResponse, error)
	CreateKeyPairsWithAddresses(ctx context.Context, req *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error)
	SignTransactionMessage(ctx context.Context, req *wallet.GetSignTransactionMessageRequest) (*wallet.GetSignTransactionMessageResponse, error)
	BuildAndSignTransaction(ctx context.Context, req *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error)
	BuildAndSignBatchTransaction(ctx context.Context, req *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error)
}
