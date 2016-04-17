package main

import (
	"fmt"
	"net"
	"encoding/json"
	"ndn/packet"
	"bytes"
	"time"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

func test() {
	fmt.Println("-----------------")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("failed to generate key", err)
	}
	publicKey := privateKey.PublicKey

	fmt.Println("public key", publicKey)
	fmt.Println("private key", privateKey)

	data := "Hello world"

	hash := sha256.New()
	msg := []byte (data)
	label := []byte("testing data")
	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, &publicKey, msg, label)
	if err != nil {
		fmt.Println("failed to encrypt test", err)
	}
	receivedText, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, cipherText, label)

	fmt.Println("data", data)
	fmt.Println("encrypted", string(cipherText))
	fmt.Println("decrypted", string(receivedText))

	fmt.Println("-----------------")
}
func main() {
	test()
	fmt.Println("NDN application demo - ping pong start")

	/* connect to proxy */
	fmt.Println("connect to proxy")
	conn, err := net.Dial("tcp", "127.0.0.1:8125")
	if err != nil {
		fmt.Println("failed to connect to proxy", err)
		panic(1)
	}
	defer conn.Close()
	fmt.Println("connected to proxy")

	/* bind data name */
	fmt.Println("bind data name")
	b := new(bytes.Buffer)
	contentName := packet.ContentName_s{}
	p := packet.ServiceProviderPacket_s{
		ContentName:contentName,
		ExpireTime:time.Now().Add(time.Second * 10),
		AllowCache:true,
		PublicKey:"",
	}
	err = json.NewEncoder(b).Encode(p)
	if (err != nil) {
		fmt.Println("failed to encode packet into json bytes")
		panic(2)
	}
	_, err = conn.Write(b.Bytes())
	if (err != nil) {
		fmt.Println("failed to send packet")
		panic(3)
	}

	/* wait for request */
	fmt.Println("wait for request")

	/* response data */
	fmt.Println("response data")

	fmt.Println("NDN application demo - ping pong end")
}