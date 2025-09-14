package solana

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

const ChainName = "Solana"

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
	return &wallet.GetChainSignMethodResponse{
		Code:       wallet.ReturnCode_SUCCESS,
		Message:    "get sign method success",
		SignMethod: "eddsa",
	}, nil
}

func (c ChainAdaptor) GetChainSchema(ctx context.Context, req *wallet.GetChainSchemaRequest) (*wallet.GetChainSchemaResponse, error) {
	ss := SolanaSchema{
		Nonce:           "",
		GasPrice:        "",
		GasTipCap:       "",
		GasFeeCap:       "",
		Gas:             0,
		ContractAddress: "",
		FromAddress:     "",
		ToAddress:       "",
		TokenId:         "",
		Value:           "",
	}
	b, err := json.Marshal(ss)
	if err != nil {
		log.Error("marshal fail", "err", err)
	}
	return &wallet.GetChainSchemaResponse{
		Code:    wallet.ReturnCode_SUCCESS,
		Message: "get Solana sign schema success",
		Schema:  string(b),
	}, nil
}

func (c ChainAdaptor) CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairsExportPublicKeyListRequest) (*wallet.CreateKeyPairsExportPublicKeyListResponse, error) {
	resp := &wallet.CreateKeyPairsExportPublicKeyListResponse{Code: wallet.ReturnCode_ERROR}
	if req.KeyNum > 10000 {
		resp.Msg = "Number must be less than 100000"
		return resp, nil
	}
	var keyList []leveldb.Key
	var retKeyList []*wallet.PublicKey
	for count := 0; count < int(req.KeyNum); count++ {
		privateKey, pubKey, compressPubkeyStr, err := c.signer.CreateKeyPair()
		if err != nil {
			if req.KeyNum > 10000 {
				resp.Msg = "create key pair fail"
				return resp, nil
			}
		}
		keyItem := leveldb.Key{
			PrivateKey: privateKey,
			Pubkey:     pubKey,
		}
		pubKeyItem := &wallet.PublicKey{
			Pubkey:         pubKey,
			CompressPubkey: compressPubkeyStr,
		}
		retKeyList = append(retKeyList, pubKeyItem)
		keyList = append(keyList, keyItem)
	}
	if ok := c.db.StoreKeys(keyList); !ok {
		resp.Msg = "create key pair fail"
		return resp, nil
	}
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Msg = "create key pair success"
	resp.PublicKeyList = retKeyList
	return resp, nil
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
