package main

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "crypto/rand"
  "crypto/rsa"
  "encoding/json"
  "fmt"
  "github.com/aabbcc1241/goutils/log"
  "io"
  "net"
  "sync"
  "time"
)

var privateKey rsa.PrivateKey

const (
  size        = 2048
  proxy_mode  = "tcp"
  proxy_addr  = "127.0.0.1:8123"
  allow_cache = true
  cache_time  = time.Second * 10
)

func init() {
  key, err := rsa.GenerateKey(rand.Reader, size)
  if err != nil {
    fmt.Println()
    panic(1)
  }
  privateKey = *key
}
func registerService(encoder json.Encoder) (err error) {

  /* bind data name */
  fmt.Println("bind data name")
  contentName := contentname.ContentName_s{
    Name:        "ping",
    ContentType: contentname.ExactMatch,
  }
  publicKey, err := packet.ToPublicKey_s(privateKey.PublicKey)
  if err != nil {
    log.Error.Println(err)
    panic(3)
  }
  out_packet := packet.ServiceProviderPacket_s{
    ContentName: contentName,
    PublicKey:   publicKey,
  }
  err = encoder.Encode(out_packet)
  if err != nil {
    fmt.Println("failed to encode packet into json bytes")
    panic(4)
  }
  fmt.Println("packet sent to proxy successfully")
  return
}
func loopForInterestPacket(wg sync.WaitGroup, in json.Decoder, out json.Encoder) (err error) {
  defer wg.Done()
  var in_packet packet.InterestPacket_s
  for err == nil {
    /* wait for request */
    fmt.Println("wait for request (incoming interest packet)")
    err = in.Decode(&in_packet)
    if err != nil {
      if err != io.EOF {
        fmt.Println("failed to decode incoming interest packet", err)
      }
      return nil
    } else {
      fmt.Println("received interest packet", in_packet)
      if in_packet.ContentName.Name == "ping" {
        onDataRequest(in_packet, out)
      }
    }
  }
  return
}
func onDataRequest(in_packet packet.InterestPacket_s, out json.Encoder) {
  /* response data */
  fmt.Println("responsing data")
  publicKey, err := packet.ToPublicKey_s(privateKey.PublicKey)
  if err != nil {
    log.Error.Println(err)
    panic(5)
  }
  out_packet := packet.DataPacket_s{
    ContentName:        in_packet.ContentName,
    SeqNum:             in_packet.SeqNum,
    AllowCache:         allow_cache && in_packet.AllowCache,
    PublisherPublicKey: publicKey,
    ContentData:        []byte("pong"),
  }
  if out_packet.AllowCache {
    out_packet.ExpireTime = time.Now().Add(cache_time)
  }
  out.Encode(out_packet)
  fmt.Println("responsed data")
}
func main() {
  //test()
  fmt.Println("NDN application demo - ping-pong server start")
  wg := sync.WaitGroup{}

  /* init connection to proxy */
  fmt.Println("connect to proxy")
  serviceConn, err := net.Dial(proxy_mode, proxy_addr)
  if err != nil {
    fmt.Println("failed to connect to proxy service socket", err)
    panic(1)
  }
  defer serviceConn.Close()
  dataConn, err := net.Dial(proxy_mode, proxy_addr)
  if err != nil {
    fmt.Println("failed to connect to proxy data socket", err)
    panic(2)
  }
  fmt.Println("connected to proxy")

  /* bind service */
  registerService(*json.NewEncoder(serviceConn))

  /* wait for request */
  wg.Add(1)
  go loopForInterestPacket(wg, *json.NewDecoder(serviceConn), *json.NewEncoder(dataConn))

  wg.Wait()
  fmt.Println("NDN application demo - ping-pong server end")
}
