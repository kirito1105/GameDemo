package tokenRSA

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func SignRsa(msg string, priKey rsa.PrivateKey) ([]byte, error) {
	//处理消息的哈希
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(msg))
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, &priKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func CheckRsa(msg string, pubKey rsa.PublicKey, signature []byte) bool {
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)
	err = rsa.VerifyPSS(&pubKey, crypto.SHA256, msgHashSum, []byte(signature), nil)

	return err == nil
}
