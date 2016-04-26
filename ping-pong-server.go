package main

import (
	"bitbucket.org/polyu-named-data-network/ndn/packet"
	"bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
	"bitbucket.org/polyu-named-data-network/ndn/packet/packettype"
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
	if err != nil {
		log.Error.Println(err)
		panic(3)
	}
	out_packet := packet.ServiceProviderPacket_s{
		ContentName: contentName,
		PublicKey:   privateKey.PublicKey,
	}
	bs, err := json.Marshal(out_packet)
	if err != nil {
		log.Error.Println("failed to marshal service provider packet", err)
		panic(6)
	}
	err = encoder.Encode(packet.GenericPacket_s{
		PacketType: packettype.ServiceProviderPacket_c,
		Payload:    bs,
	})
	if err != nil {
		fmt.Println("failed to encode packet into json bytes")
		panic(4)
	}
	fmt.Println("packet sent to proxy successfully")
	return
}

func onDataRequest(in_packet packet.InterestPacket_s, out json.Encoder) {
	/* response data */
	log.Info.Println("responsing data", in_packet.ContentName.Name)
	out_packet := packet.DataPacket_s{
		ContentName:        in_packet.ContentName,
		SeqNum:             in_packet.SeqNum,
		AllowCache:         allow_cache && in_packet.AllowCache,
		PublisherPublicKey: privateKey.PublicKey,
		ContentData:        []byte("pong"),
	}
	if out_packet.AllowCache {
		out_packet.ExpireTime = time.Now().Add(cache_time)
	}
	bs, err := json.Marshal(out_packet)
	if err != nil {
		log.Error.Println("failed to marshal data packet", err)
		panic(7)
	}
	out.Encode(packet.GenericPacket_s{
		PacketType: packettype.DataPacket_c,
		Payload:    bs,
	})
	log.Info.Println("responsed data")
}
func init() {
	log.Init(true, true, true, log.DefaultCommFlag)
}
func main() {
	log.Info.Println("NDN application demo - ping-pong server start")
	wg := sync.WaitGroup{}

	/* init connection to proxy */
	log.Info.Println("connect to proxy")
	conn, err := net.Dial(proxy_mode, proxy_addr)
	if err != nil {
		log.Error.Println("failed to connect to proxy", proxy_mode, proxy_addr, err)
		panic(1)
	}
	defer conn.Close()
	log.Info.Println("connected to proxy")

	/* bind service */
	encoder := json.NewEncoder(conn)
	registerService(*encoder)

	/* wait for request */
	decoder := json.NewDecoder(conn)
	for err == nil {
		var in_packet = packet.GenericPacket_s{}
		err = decoder.Decode(&in_packet)
		if err != nil {
			if err != io.EOF {
				log.Error.Println("failed to parse incoming packet")
			}
		} else {
			if in_packet.PacketType == packettype.InterestPacket_c {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var p packet.InterestPacket_s
					err = json.Unmarshal(in_packet.Payload, &p)
					if err != nil {
						log.Error.Println("failed to parse interest packet", err, in_packet)
					}
					if p.ContentName.Name == "ping" {
						onDataRequest(p, *encoder)
					} else {
						log.Info.Println("received unexpected interest for", p.ContentName.Name)
						p := packet.InterestReturnPacket_s{}
						if gp, err := p.ToGenericPacket(); err != nil {
							log.Error.Println("failed to marshal inerest return packet", err)
						} else {
							err = encoder.Encode(gp)
							if err != nil {
								log.Error.Println("failed to send interest return packet")
							}
						}
					}
				}()
			} else {
				log.Error.Println("unexpected packet", in_packet)
			}
		}
	}

	wg.Wait()
	log.Info.Println("NDN application demo - ping-pong server end")
}
