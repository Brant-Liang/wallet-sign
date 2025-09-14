package ethereum

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Brant-Liang/wallet-sign/chain"
	"github.com/Brant-Liang/wallet-sign/config"
	wallet "github.com/Brant-Liang/wallet-sign/gen/go"
	"github.com/Brant-Liang/wallet-sign/hsm"
	"github.com/Brant-Liang/wallet-sign/leveldb"
	"github.com/Brant-Liang/wallet-sign/ssm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"
)

const (
	ChainName            = "Ethereum"
	maxCreateKeyPairsNum = 10_000
)

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
		SignMethod: ssm.ECDSA,
	}, nil
}

func (c ChainAdaptor) GetChainSchema(ctx context.Context, req *wallet.GetChainSchemaRequest) (*wallet.GetChainSchemaResponse, error) {
	es := EthereumSchema{
		RequestId: "0",
		DynamicFeeTx: Eip1559DynamicFeeTx{
			ChainId:              "",
			Nonce:                0,
			FromAddress:          common.Address{}.String(),
			ToAddress:            common.Address{}.String(),
			GasLimit:             0,
			MaxFeePerGas:         "0",
			MaxPriorityFeePerGas: "0",
			Amount:               "0",
			ContractAddress:      "",
		},
		ClassicFeeTx: LegacyFeeTx{
			ChainId:         "0",
			Nonce:           0,
			FromAddress:     common.Address{}.String(),
			ToAddress:       common.Address{}.String(),
			GasLimit:        0,
			GasPrice:        0,
			Amount:          "0",
			ContractAddress: "",
		},
	}
	b, err := json.Marshal(es)
	if err != nil {
		log.Error("marshal fail", "err", err)
	}
	return &wallet.GetChainSchemaResponse{
		Code:    wallet.ReturnCode_SUCCESS,
		Message: "get ethereum sign schema success",
		Schema:  string(b),
	}, nil
}

func (c ChainAdaptor) CreateKeyPairsExportPublicKeyList(ctx context.Context, req *wallet.CreateKeyPairsExportPublicKeyListRequest) (*wallet.CreateKeyPairsExportPublicKeyListResponse, error) {
	resp := &wallet.CreateKeyPairsExportPublicKeyListResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	if req.KeyNum <= 0 {
		resp.Msg = "key number must be greater than 0"
		return resp, nil
	}
	if req.KeyNum > maxCreateKeyPairsNum {
		resp.Msg = fmt.Sprintf("number must be <= %d", maxCreateKeyPairsNum)
		return resp, nil
	}
	if c.signer == nil {
		return nil, errors.New("signer not initialized")
	}
	if c.db == nil {
		return nil, errors.New("db not initialized")
	}
	var keyList []leveldb.Key
	var retKeyList []*wallet.PublicKey
	for counter := 0; counter < int(req.KeyNum); counter++ {
		priKeyStr, pubKeyStr, compressPubkeyStr, err := c.signer.CreateKeyPair()
		if err != nil {
			resp.Msg = "create key pairs fail"
			return resp, nil
		}
		keyItem := leveldb.Key{
			PrivateKey: priKeyStr,
			Pubkey:     pubKeyStr,
		}
		pukItem := &wallet.PublicKey{
			CompressPubkey: compressPubkeyStr,
			Pubkey:         pubKeyStr,
		}
		retKeyList = append(retKeyList, pukItem)
		keyList = append(keyList, keyItem)
	}
	if ok := c.db.StoreKeys(keyList); !ok {
		log.Error("store keys fail", "isOk", ok)
		return nil, errors.New("store keys fail")
	}
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Msg = "create keys success"
	resp.PublicKeyList = retKeyList
	return resp, nil
}

func (c ChainAdaptor) CreateKeyPairsWithAddresses(ctx context.Context, req *wallet.CreateKeyPairsWithAddressesRequest) (*wallet.CreateKeyPairsWithAddressesResponse, error) {
	resp := &wallet.CreateKeyPairsWithAddressesResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	if req.KeyNum > 10000 {
		resp.Message = "Number must be less than 100000"
		return resp, nil
	}
	var keyList []leveldb.Key
	var retKeyWithAddressList []*wallet.ExportPublicKeyWithAddress
	for counter := 0; counter < int(req.KeyNum); counter++ {
		priKeyStr, pubKeyStr, compressPubkeyStr, err := c.signer.CreateKeyPair()
		if err != nil {
			resp.Message = "create key pairs fail"
			return resp, nil
		}
		keyItem := leveldb.Key{
			PrivateKey: priKeyStr,
			Pubkey:     pubKeyStr,
		}
		publicKeyBytes, err := hex.DecodeString(pubKeyStr)
		pukAddressItem := &wallet.ExportPublicKeyWithAddress{
			CompressPublicKey: compressPubkeyStr,
			PublicKey:         pubKeyStr,
			Address:           common.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:]).String(),
		}
		retKeyWithAddressList = append(retKeyWithAddressList, pukAddressItem)
		keyList = append(keyList, keyItem)
	}
	if ok := c.db.StoreKeys(keyList); !ok {
		log.Error("store keys fail", "isOk", ok)
		return nil, errors.New("store keys fail")
	}
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Message = "create key pairs success"
	resp.PublicKeyAddresses = retKeyWithAddressList
	return resp, nil
}

func (c ChainAdaptor) SignTransactionMessage(ctx context.Context, req *wallet.GetSignTransactionMessageRequest) (*wallet.GetSignTransactionMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildAndSignBatchTransaction(ctx context.Context, req *wallet.BuildAndSignBatchTransactionRequest) (*wallet.BuildAndSignBatchTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildAndSignTransaction(ctx context.Context, req *wallet.BuildAndSignTransactionRequest) (*wallet.BuildAndSignTransactionResponse, error) {
	resp := &wallet.BuildAndSignTransactionResponse{Code: wallet.ReturnCode_ERROR}

	// 1) 解析 & 构造 tx
	dFeeTx, _, err := c.buildDynamicFeeTx(req.TxBase64Body)
	if err != nil {
		return nil, err
	}
	// 2) 待签名hash (digest)
	digest := CreateEip1559UnSignTx(dFeeTx, dFeeTx.ChainID) // 返回 common.Hash

	rawTx := digest.Hex() // digest: 对 TxData 规范化编码 + keccak256，结果 32字节

	// 取私钥
	privKey, ok := c.db.GetPrivKey(req.PublicKey)
	if !ok {
		log.Error("private key not found for public key")
		resp.Message = "private key not found"
		return resp, nil
	}

	// 4) 签名（约定 SignMessage 接收 32字节哈希的 hex 字符串，不带0x）
	//    如果需要，去掉 0x：hex.EncodeToString(digest[:])
	signature, err3 := c.signer.SignMessage(privKey, rawTx)
	if err3 != nil {
		log.Error("sign transaction fail", "err", err3)
		resp.Message = "sign transaction fail"
		return resp, nil
	}

	inputSignatureByteList, err := hex.DecodeString(signature)
	if err != nil {
		log.Error("decode signature failed", "err", err)
		resp.Message = "decode signature failed"
		return resp, nil
	}

	eip1559Signer, signedTx, signAndHandledTx, txHash, err := CreateEip1559SignedTx(dFeeTx, inputSignatureByteList, dFeeTx.ChainID)
	if err != nil {
		log.Error("create signed tx fail", "err", err)
		resp.Message = "create signed tx fail"
		return resp, nil
	}
	log.Info("sign transaction success",
		"eip1559Signer", eip1559Signer,
		"signedTx", signedTx,
		"signAndHandledTx", signAndHandledTx,
		"txHash", txHash,
	)
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Message = "sign whole transaction success"
	resp.SignedTx = signAndHandledTx
	resp.TxHash = txHash
	resp.TxMessageHash = rawTx
	return resp, nil
}

func (c ChainAdaptor) buildDynamicFeeTx(base64Tx string) (*types.DynamicFeeTx, *Eip1559DynamicFeeTx, error) {
	// 1. Decode base64 string
	txReqJsonByte, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, nil, err
	}

	// 2. Unmarshal JSON to struct
	var dynamicFeeTx Eip1559DynamicFeeTx
	if err := json.Unmarshal(txReqJsonByte, &dynamicFeeTx); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, nil, err
	}

	// 3. Convert string values to big.Int
	chainID := new(big.Int)
	maxPriorityFeePerGas := new(big.Int)
	maxFeePerGas := new(big.Int)
	amount := new(big.Int)

	log.Info("Dynamic fee tx",
		"ChainId", dynamicFeeTx.ChainId,
		"MaxPriorityFeePerGas", dynamicFeeTx.MaxPriorityFeePerGas,
		"MaxFeePerGas", dynamicFeeTx.MaxFeePerGas,
		"Amount", dynamicFeeTx.Amount,
	)

	if _, ok := chainID.SetString(dynamicFeeTx.ChainId, 10); !ok {
		return nil, nil, fmt.Errorf("invalid chain ID: %s", dynamicFeeTx.ChainId)
	}
	if _, ok := maxPriorityFeePerGas.SetString(dynamicFeeTx.MaxPriorityFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max priority fee: %s", dynamicFeeTx.MaxPriorityFeePerGas)
	}
	if _, ok := maxFeePerGas.SetString(dynamicFeeTx.MaxFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max fee: %s", dynamicFeeTx.MaxFeePerGas)
	}
	if _, ok := amount.SetString(dynamicFeeTx.Amount, 10); !ok {
		return nil, nil, fmt.Errorf("invalid amount: %s", dynamicFeeTx.Amount)
	}

	// 4. Handle addresses and data
	toAddress := common.HexToAddress(dynamicFeeTx.ToAddress)
	var finalToAddress common.Address
	var finalAmount *big.Int
	var buildData []byte
	log.Info("contract address check",
		"contractAddress", dynamicFeeTx.ContractAddress,
		"isEthTransfer", isEthTransfer(&dynamicFeeTx),
	)

	// 5. Handle contract interaction vs direct transfer
	if isEthTransfer(&dynamicFeeTx) {
		log.Info("native token transfer")
		finalToAddress = toAddress
		finalAmount = amount
	} else {
		log.Info("erc20 token transfer")
		contractAddress := common.HexToAddress(dynamicFeeTx.ContractAddress)
		buildData = BuildErc20Data(toAddress, amount)
		finalToAddress = contractAddress
		finalAmount = big.NewInt(0)
	}

	// 6. Create dynamic fee transaction
	dFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     dynamicFeeTx.Nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       dynamicFeeTx.GasLimit,
		To:        &finalToAddress,
		Value:     finalAmount,
		Data:      buildData,
	}

	return dFeeTx, &dynamicFeeTx, nil
}

func isEthTransfer(tx *Eip1559DynamicFeeTx) bool {
	if tx.ContractAddress == "" || strings.ToLower(tx.ContractAddress) == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" {
		return true
	}
	return false
}
