package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {

		panic(err)
	}

	// The public key is a part of the *rsa.PrivateKey struct
	//publicKey := privateKey.PublicKey
	str, _ := json.Marshal(privateKey)
	file, err := os.Create("rcenterServer/key.private.pem")

	_, err = file.Write(str)
	if err != nil {
		fmt.Println(err)
		return
	}

	publicKey := privateKey.PublicKey
	str2, _ := json.Marshal(publicKey)
	file, err = os.Create("roomServer/key.public.pem")
	_, err = file.Write(str2)
	if err != nil {
		return
	}
	fmt.Println(privateKey)
	fmt.Println(publicKey)
}
