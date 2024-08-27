package main

import (
	"flag"
	"net"

	"github.com/liyouxina/tunnel/pkg/logger"
)

var tunnelPort = flag.String("tunnelPort", "8091", "tunnelPort")
var waiCon net.Conn
var neiCon net.Conn
var log = logger.Logger

func init() {
	flag.Parse()
}

func main() {
	var err error
	listener, err := net.Listen("tcp", ":"+*tunnelPort)
	if err != nil {
		log.Fatalf("start tunnel server failed %v", err)
		return
	}
	log.Infof("start tunnel server success %v", *tunnelPort)
	for {
		neiCon, err = listener.Accept()
		if err != nil {
			log.Warnf("Failed to accept tunnel connection %v", err)
			continue
		}

		go func() {
			for {
				buffer := make([]byte, 1024)
				n, err := neiCon.Read(buffer)
				if err != nil {
					log.Warnf("Failed to read nei tunnel connection %v", err)
				}
				_, err = waiCon.Write(buffer[:n])
				if err != nil {
					log.Warnf("Failed to write wai tunnel connection %v", err)
				}
			}

		}()
	}

}

func server() {
	var err error
	listener, err := net.Listen("tcp", ":10001")
	if err != nil {
		log.Fatalf("start tunnel server failed %v", err)
		return
	}
	log.Infof("start tunnel server success %v", *tunnelPort)
	for {
		waiCon, err = listener.Accept()
		if err != nil {
			log.Warnf("Failed to accept tunnel connection %v", err)
			continue
		}

		go func() {
			for {
				buffer := make([]byte, 1024)
				n, err := waiCon.Read(buffer)
				if err != nil {
					log.Warnf("Failed to read wai tunnel connection %v", err)
				}
				_, err = neiCon.Write(buffer[:n])
				if err != nil {
					log.Warnf("Failed to write nei tunnel connection %v", err)
				}
			}

		}()
	}
}
