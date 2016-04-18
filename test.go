package main

import (
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "encoding/json"
  "fmt"
  "io"
  "ndn/packet"
  "ndn/packet/contentname"
  "net"
)

func test() {
  fmt.Println("-----------------")

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

var privateKey rsa.PrivateKey

const size = 2048

func init() {
  key, err := rsa.GenerateKey(rand.Reader, size)
  if err != nil {
    fmt.Println()
    panic(1)
  }
  privateKey = *key
}
