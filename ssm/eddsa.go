package ssm

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
)

type EdDSAKeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func CreateEdDSAKeyPair() (string, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Error("create key pair fail:", "err", err)
		return EmptyHexString, EmptyHexString, nil
	}
	return hex.EncodeToString(privateKey), hex.EncodeToString(publicKey), nil
}

func ParseEdDSAPublicKey(hexKey string) (ed25519.PublicKey, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid public key hex: %w", err)
	}
	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}
	return ed25519.PublicKey(keyBytes), nil
}

func SignEdDSAMessage(priKey string, txMsg string) (string, error) {
	privateKey, _ := hex.DecodeString(priKey)
	txMsgByte, _ := hex.DecodeString(txMsg)
	signMsg := ed25519.Sign(privateKey, txMsgByte)

	return hex.EncodeToString(signMsg), nil
}

func VerifyEdDSASign(pubKey, msgHash, sig string) bool {
	publicKeyByte, _ := hex.DecodeString(pubKey)
	msgHashByte, _ := hex.DecodeString(msgHash)
	signature, _ := hex.DecodeString(sig)
	return ed25519.Verify(publicKeyByte, msgHashByte, signature)
}
