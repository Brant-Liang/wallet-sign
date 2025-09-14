package ssm

import (
	"errors"
)

const (
	EmptyHexString = "0x00"
)

// CryptoType Define a custom type for cryptographic algorithm types
type CryptoType = string

// Define constants for the supported cryptographic types
const (
	ECDSA CryptoType = "ecdsa"
	EDDSA CryptoType = "eddsa"
)

func ParseTransactionType(s string) (CryptoType, error) {
	switch s {
	case ECDSA:
		return ECDSA, nil
	case EDDSA:
		return EDDSA, nil
	default:
		return "", errors.New("unknown transaction type")
	}
}

func GetSignerByType(cryptoType CryptoType) (Signer, error) {
	switch cryptoType {
	case ECDSA:
		return NewEcdsaSigner(), nil
	case EDDSA:
		return NewEdDSASigner(), nil
	default:
		return nil, errors.New("unsupported crypto type: ")
	}
}
