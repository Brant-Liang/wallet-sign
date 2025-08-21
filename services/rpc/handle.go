package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/Brant-Liang/wallet-sign/gen/go"
	"github.com/Brant-Liang/wallet-sign/leveldb"
	"github.com/Brant-Liang/wallet-sign/ssm"
)

const BearerToken = "BearerToken"

func (s *RpcServer) GetSupportSignType(ctx context.Context, in *wallet.GetSupportSignWayRequest) (*wallet.GetSupportSignWayResponse, error) {
	var signWay []*wallet.SignWay

	if in.ConsumerToken != BearerToken {
		return &wallet.GetSupportSignWayResponse{
			Code:    wallet.ReturnCode_ERROR,
			Msg:     "get sign way fail",
			SignWay: signWay,
		}, nil
	}
	signWay = append(signWay, &wallet.SignWay{Schema: "ecdsa"})
	signWay = append(signWay, &wallet.SignWay{Schema: "eddsa"})
	return &wallet.GetSupportSignWayResponse{
		Code:    wallet.ReturnCode_SUCCESS,
		Msg:     "get sign way success",
		SignWay: signWay,
	}, nil
}

func (s *RpcServer) CreateKeyPairsExportPublicKeyList(ctx context.Context, in *wallet.ExportPublicKeyRequest) (*wallet.ExportPublicKeyResponse, error) {
	resp := &wallet.ExportPublicKeyResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	if in.ConsumerToken != BearerToken {
		resp.Msg = "bearer token fail"
		return resp, nil
	}
	cryptoType, err := ssm.ParseTransactionType(in.Type)
	if err != nil {
		resp.Msg = "input sign type error"
		return resp, nil
	}
	if in.Number > 10000 {
		resp.Msg = "Number must be less than 100000"
		return resp, nil
	}

	var keyList []leveldb.Key
	var pubKeyList []*wallet.PublicKey

	for counter := 0; counter < int(in.Number); counter++ {
		var priKeyStr, pubKeyStr, compressPubkeyStr string
		var err error
		switch cryptoType {
		case ssm.ECDSA:
			priKeyStr, pubKeyStr, compressPubkeyStr, err = ssm.CreateECDSAKeyPair()
		case ssm.EDDSA:
			priKeyStr, pubKeyStr, err = ssm.CreateEdDSAKeyPair()
			compressPubkeyStr = pubKeyStr
		default:
			return nil, errors.New("unsupported key type")
		}
		if err != nil {
			log.Error("create key pair fail", "err", err)
			return nil, err
		}
		keyItem := leveldb.Key{
			PrivateKey: priKeyStr,
			Pubkey:     pubKeyStr,
		}
		pukItem := &wallet.PublicKey{
			CompressPubkey: compressPubkeyStr,
			Pubkey:         pubKeyStr,
		}
		pubKeyList = append(pubKeyList, pukItem)
		keyList = append(keyList, keyItem)
	}
	isOk := s.db.StoreKeys(keyList)
	if !isOk {
		log.Error("store keys fail", "isOk", isOk)
		return nil, errors.New("store keys fail")
	}
	resp.Code = wallet.ReturnCode_SUCCESS
	resp.Msg = "create key pairs success"
	resp.PublicKey = pubKeyList
	return resp, nil
}

func (s *RpcServer) SignMessageSignature(ctx context.Context, in *wallet.SignTxMessageRequest) (*wallet.SignTxMessageResponse, error) {
	resp := &wallet.SignTxMessageResponse{
		Code: wallet.ReturnCode_ERROR,
	}
	if in.ConsumerToken != BearerToken {
		resp.Msg = "Bearer Token Error"
		return resp, nil
	}
	cryptotype, err := ssm.ParseTransactionType(in.Type)
	if err != nil {
		resp.Msg = "input type error"
		return resp, nil
	}
	privatekey, isOk := s.db.GetPrivKey(in.PublicKey)
	if !isOk {
		resp.Msg = "private key error"
		return resp, nil
	}
	var signature string
	var err2 error
	switch cryptotype {
	case ssm.ECDSA:
		signature, err2 = ssm.SignECDSAMessage(privatekey, in.MessageHash)
	case ssm.EDDSA:
		signature, err2 = ssm.SignEdDSAMessage(privatekey, in.MessageHash)
	default:
		return nil, errors.New("unknown sign type")
	}
	if err2 != nil {
		return nil, err2
	}
	resp.Signature = signature
	resp.Msg = "signature success"
	resp.Hash = in.MessageHash
	resp.Code = wallet.ReturnCode_SUCCESS
	return resp, nil
}

func (s *RpcServer) GetSignBatchMessageSignature(ctx context.Context, in *wallet.SignBatchMessageSignatureRequest) (*wallet.SignBatchMessageSignatureResponse, error) {
	resp := &wallet.SignBatchMessageSignatureResponse{
		Code: wallet.ReturnCode_SUCCESS,
	}
	var msList []*wallet.MessageSignature
	for _, msghash := range in.MsgHashList {
		cryptoType, err := ssm.ParseTransactionType(msghash.SignType)
		if err != nil {
			log.Error("parse transaction error", "err", msghash.TxMessageHash)
			continue
		}
		privateKey, isOk := s.db.GetPrivKey(msghash.PublicKey)
		if !isOk {
			log.Error("get private key error", "err", err)
			continue
		}
		var signature string
		var err2 error
		switch cryptoType {
		case ssm.ECDSA:
			signature, err2 = ssm.SignECDSAMessage(privateKey, msghash.TxMessageHash)
		case ssm.EDDSA:
			signature, err2 = ssm.SignEdDSAMessage(privateKey, msghash.TxMessageHash)
		default:
			return nil, errors.New("unknown sign type")
		}
		if err2 != nil {
			log.Error("sign batch message signature error", "err", err2)
			continue
		}
		msList = append(msList, &wallet.MessageSignature{
			Signature:     signature,
			TxMessageHash: msghash.TxMessageHash,
		})
	}
	resp.SuccessMsgSigList = msList
	resp.Msg = "sign batch message signature success"
	resp.Code = wallet.ReturnCode_SUCCESS
	return resp, nil
}
