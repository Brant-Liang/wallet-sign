package ethereum

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

// ERC-20 的 transfer(address,uint256)
func BuildErc20Data(to common.Address, amount *big.Int) []byte {
	if amount == nil || amount.Sign() < 0 {
		return nil
	}
	sig := []byte("transfer(address,uint256)")
	methodID := crypto.Keccak256(sig)[:4]
	addr := common.LeftPadBytes(to.Bytes(), 32)
	amt := common.LeftPadBytes(amount.Bytes(), 32)

	data := make([]byte, 0, 4+32+32)
	data = append(data, methodID...)
	data = append(data, addr...)
	data = append(data, amt...)
	return data
}

func BuildErc721Data(from, to common.Address, tokenId *big.Int) []byte {
	if tokenId == nil || tokenId.Sign() < 0 {
		panic("invalid tokenId")
	}
	sig := []byte("safeTransferFrom(address,address,uint256)")
	methodID := crypto.Keccak256(sig)[:4]

	f := common.LeftPadBytes(from.Bytes(), 32)
	t := common.LeftPadBytes(to.Bytes(), 32)
	id := common.LeftPadBytes(tokenId.Bytes(), 32)

	data := make([]byte, 0, 4+32+32+32)
	data = append(data, methodID...)
	data = append(data, f...)
	data = append(data, t...)
	data = append(data, id...)
	return data
}

// 生成 未签名交易的哈希（digest），后续要拿去做签名的那个消息哈希
// Legacy 类型的交易数据（旧格式，包含 Nonce, GasPrice, GasLimit, To, Value, Data）。
// chainId → 用于 EIP-155 replay protection。
func CreateLegacyUnSignTx(txData *types.LegacyTx, chainId *big.Int) common.Hash {
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)
	return signer.Hash(tx)
}

// txData → EIP-1559 类型的交易数据（DynamicFeeTx，包括 MaxFeePerGas, MaxPriorityFeePerGas 等字段）。
// chainId → 指定的链 ID。
func CreateEip1559UnSignTx(txData *types.DynamicFeeTx, chainId *big.Int) common.Hash {
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)
	return signer.Hash(tx)
}

func CreateLegacySignedTx(txData *types.LegacyTx, sig []byte, chainId *big.Int) (rawHex string, txHash string, err error) {
	if len(sig) != 65 {
		return "", "", errors.New("invalid signature length")
	}
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)

	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", "", errors.Wrap(err, "with signature")
	}

	enc, err2 := rlp.EncodeToBytes(signedTx) // or signedTx.MarshalBinary()
	if err2 != nil {
		return "", "", errors.Wrap(err2, "encode rlp")
	}

	return "0x" + hex.EncodeToString(enc), signedTx.Hash().String(), nil
}

func CreateEip1559SignedTx(txData *types.DynamicFeeTx, sig []byte, chainId *big.Int) (types.Signer, *types.Transaction, string, string, error) {
	// r(32)||s(32)||v(1)
	if len(sig) != 65 {
		return nil, nil, "", "", errors.New("invalid signature length")
	}
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)

	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return nil, nil, "", "", errors.Wrap(err, "with signature")
	}

	enc, err2 := signedTx.MarshalBinary() // or signedTx.MarshalBinary()
	if err2 != nil {
		return nil, nil, "", "", errors.Wrap(err2, "encode rlp")
	}

	rawHex := "0x" + hex.EncodeToString(enc)
	return signer, signedTx, rawHex, signedTx.Hash().String(), nil
}
