package main

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "crypto/rand"
  "crypto/rsa"
  "encoding/json"
  "fmt"
  "io"
  "net"
)

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
func main() {
  //test()
  fmt.Println("NDN application demo - ping-pong server start")

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
  contentName := packet.ContentName_s{
    Name: "ping",
    Type: contentname.ExactMatch,
  }
  out_packet := packet.ServiceProviderPacket_s{
    ContentName: contentName,
    PublicKey:   privateKey.PublicKey,
  }
  err = json.NewEncoder(conn).Encode(out_packet)
  if err != nil {
    fmt.Println("failed to encode packet into json bytes")
    panic(2)
  }
  fmt.Println("packet sent to proxy successfully")

  decoder := json.NewDecoder(conn)
  var in_packet packet.InterestPacket_s
  for err == nil {
    fmt.Println("wait for incoming interest packet")
    err = decoder.Decode(&in_packet)
    if err != nil {
      if err != io.EOF {
        fmt.Println("failed to decode incoming interest packet")
      }
    } else {
      fmt.Println("received interest packet", in_packet)
    }
  }
  /* wait for request */
  fmt.Println("wait for request")

  /* response data */
  fmt.Println("response data")

  fmt.Println("NDN application demo - ping-pong server end")
}
