package ssm

type Signer interface {
	CreateKeyPair() (privateKey string, publicKey string, compressPubkeyStr string, err error)
	SignMessage(privateKey string, txMsg string) (string, error)
	VerifyMessage(publicKey string, txMsg string, signature string) (bool, error)
}
