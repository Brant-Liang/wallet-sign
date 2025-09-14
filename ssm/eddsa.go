package ssm

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/log"
)

type EdDSASinger struct{}

func NewEdDSASigner() *EdDSASinger {
	return &EdDSASinger{}
}

func (e *EdDSASinger) CreateKeyPair() (string, string, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Error("create key pair fail:", "err", err)
		return EmptyHexString, EmptyHexString, "", nil
	}
	return hex.EncodeToString(privateKey), hex.EncodeToString(publicKey), hex.EncodeToString(publicKey), nil
}

func (e *EdDSASinger) SignMessage(privateKey string, txMsg string) (string, error) {
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Error("decode ed25519 private key fail", "err", err)
		return EmptyHexString, err
	}
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return EmptyHexString, errors.New("invalid ed25519 private key length")
	}

	txMsgBytes, err := hex.DecodeString(txMsg)
	if err != nil {
		log.Error("decode tx message fail", "err", err)
		return EmptyHexString, err
	}

	signature := ed25519.Sign(privateKeyBytes, txMsgBytes)
	return hex.EncodeToString(signature), nil
}

func (e *EdDSASinger) VerifyMessage(publicKey string, txMsg string, signature string) (bool, error) {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Error("decode ed25519 public key fail", "err", err)
		return false, err
	}
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return false, errors.New("invalid ed25519 public key length")
	}

	txHashBytes, err := hex.DecodeString(txMsg)
	if err != nil {
		log.Error("decode tx hash fail", "err", err)
		return false, err
	}

	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		log.Error("decode signature fail", "err", err)
		return false, err
	}
	if len(sigBytes) != ed25519.SignatureSize {
		return false, errors.New("invalid ed25519 signature length")
	}

	ok := ed25519.Verify(publicKeyBytes, txHashBytes, sigBytes)
	return ok, nil
}
