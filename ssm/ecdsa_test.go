package ssm

import (
	"fmt"
	"testing"
)

func TestCreateECDSAKeyPair(t *testing.T) {
	var singer Signer = NewEcdsaSigner()
	privKey, pubKey, cpubKey, _ := singer.CreateKeyPair()
	fmt.Println("privKey=", privKey)
	fmt.Println("pubKey=", pubKey)
	fmt.Println("cpubKey=", cpubKey)
	// priKeyStr	私钥（hex 编码）
	// pubKeyStr	公钥（未压缩，hex 编码）
	// compressPubkeyStr	公钥（压缩格式，hex 编码）
	// error	错误对象（生成失败时返回）
}

func TestSignMessage(t *testing.T) {
	// 0x35096AD62E57e86032a3Bb35aDaCF2240d55421D
	var signer Signer = NewEcdsaSigner()
	privKey := "fb26155c1ff94bb97692793d1197d9c6c8091f25f8c8ac703f92695d32c5194b"
	message := "0x3e4f9a460233ec33862da1ac3dabf5b32db01400fba166cdec40ad6dc735b4ab"
	signature, err := signer.SignMessage(privKey, message)
	if err != nil {
		fmt.Println("sign tx fail")
	}
	fmt.Println("Signature: ", signature)
}

func TestVerifyEcdsaSignature(t *testing.T) {
	var signer Signer = NewEcdsaSigner()
	CompressedPubKey := "028846b3ce4376e8d58c83c1c6420a784caa675d7f26c496f499585d09891af8fc"
	txMsg := "3e4f9a460233ec33862da1ac3dabf5b32db01400fba166cdec40ad6dc735b4ab"
	signature := "f8c9ab615ffd81f74d9db8765e25ce260ba3b4da1c6af2a52dedc697dcff833b6cfe576a1b6b7106a6880d8057639d4b87a67001c69594df29d928d6048912f900"

	isValid, err := signer.VerifyMessage(CompressedPubKey, txMsg, signature)
	if err != nil {
		t.Error("Failed to verify signature:", err)
	}

	if !isValid {
		t.Error("Signature is invalid")
	}
}
