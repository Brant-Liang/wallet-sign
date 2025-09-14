package ssm

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type EcdsaSigner struct{}

func NewEcdsaSigner() *EcdsaSigner {
	return &EcdsaSigner{}
}

func (s *EcdsaSigner) CreateKeyPair() (string, string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Error("generate key fail", "err", err)
		return EmptyHexString, EmptyHexString, EmptyHexString, err
	}
	// hex EncodeToString 转为十六进制字符串
	priKeyStr := hex.EncodeToString(crypto.FromECDSA(privateKey))
	pubKeyStr := hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey))
	compressPubkeyStr := hex.EncodeToString(crypto.CompressPubkey(&privateKey.PublicKey))
	return priKeyStr, pubKeyStr, compressPubkeyStr, err
}

func (s *EcdsaSigner) SignMessage(privateKey string, txMsg string) (string, error) {
	hash := common.HexToHash(txMsg)
	privByte, err := hex.DecodeString(privateKey) //私钥的 hex 编码解码成 byte 数组，通常是 32 字节。
	if err != nil {
		log.Error("decode private key fail", "err", err)
		return EmptyHexString, err
	}
	privKeyEcdsa, err := crypto.ToECDSA(privByte)
	// 将 byte 形式私钥恢复成 *ecdsa.PrivateKey
	// crypto.ToECDSA() 会自动验证私钥是否合法（在 secp256k1 上）
	if err != nil {
		log.Error("Byte private key to ecdsa key fail", "err", err)
		return EmptyHexString, err
	}
	signatureByte, err := crypto.Sign(hash[:], privKeyEcdsa)
	if err != nil {
		log.Error("sign transaction fail", "err", err)
		return EmptyHexString, err
	}
	return hex.EncodeToString(signatureByte), nil
}

func (s *EcdsaSigner) VerifyMessage(publicKey string, txMsg string, signature string) (bool, error) {
	// Convert public key from hexadecimal to bytes
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Error("Error converting public key to bytes", err)
		return false, err
	}

	// Convert transaction string from hexadecimal to bytes
	txMsgBytes, err := hex.DecodeString(txMsg)
	if err != nil {
		log.Error("Error converting transaction hash to bytes", err)
		return false, err
	}

	// Convert signature from hexadecimal to bytes
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Error("Error converting signature to bytes", err)
		return false, err
	}

	// Verify the transaction signature using the public key
	return crypto.VerifySignature(pubKeyBytes, txMsgBytes, sigBytes[:64]), nil
}
