package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func key_test() {
	fmt.Println("-----------------")

	const size = 2048

	var privateKey rsa.PrivateKey
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		fmt.Println()
		panic(1)
	}
	privateKey = *key

	//privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	//if err != nil {
	//  fmt.Println("failed to generate key", err)
	//}

	publicKey := privateKey.PublicKey

	fmt.Println("public key", publicKey)
	fmt.Println("private key", privateKey)

	data := "Hello world"

	hash := sha256.New()
	msg := []byte(data)
	label := []byte("testing data")
	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, &publicKey, msg, label)
	if err != nil {
		fmt.Println("failed to encrypt test", err)
	}
	receivedText, err := rsa.DecryptOAEP(hash, rand.Reader, &privateKey, cipherText, label)

	fmt.Println("data", data)
	fmt.Println("encrypted", string(cipherText))
	fmt.Println("decrypted", string(receivedText))

	fmt.Println("-----------------")
}
func byte_test() {
	ContentData := []byte("pong")
	fmt.Println("the content raw", ContentData)
	fmt.Println("the content string", string(ContentData))
}
func main() {
	byte_test()
}
