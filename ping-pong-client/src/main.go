package main

import (
  "encoding/json"
  "fmt"
  "ndn/packet"
  "net"
  "sync"
)

func main() {
  fmt.Println("NDN application demo - ping-pong client start")

  wg := sync.WaitGroup{}

  /* establish interest connection */
  interestConn, err := net.Dial("tcp", "127.0.0.1:8123")
  if err != nil {
    fmt.Println("failed to connect to proxy interest service", err)
    panic(1)
  }

  fmt.Println("preparing interest packet")
  out_packet := packet.InterestPacket_s{}
  err = json.NewEncoder(interestConn).Encode(out_packet)
  if err != nil {
    fmt.Println("failed to encode interest packet", err)
    panic(2)
  }
  fmt.Println("sent interest packet")

  /* prepare interest packet (request) */
  wg.Add(1)
  go func() {
    defer wg.Done()
    var in_packet packet.InterestReturnPacket_s
    fmt.Println("wait for interestReturn packet")
    json.NewDecoder(interestConn).Decode(&in_packet)
    fmt.Println("received interestReturn pcaket", in_packet)
  }()

  /* wait for interest return (NAK) */

  /* establish data connection */
  dataConn, err := net.Dial("tcp", "127.0.0.1:8124")
  if err != nil {
    fmt.Println("failed to connect to proxy data service", err)
    panic(3)
  }
  wg.Add(1)
  go func() {
    defer wg.Done()
    var in_packet packet.DataPacket_s
    fmt.Println("wait for data packet")
    json.NewDecoder(dataConn).Decode(&in_packet)
    fmt.Println("received data packet", in_packet)
  }()

  /* wait for data packet (response)*/

  wg.Wait()
  fmt.Println("NDN application demo - ping-pong client end")
}
