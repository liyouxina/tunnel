package main

import (
	"flag"
	"net"

	"github.com/liyouxina/tunnel/pkg/logger"
)

var waiCon net.Conn
var neiCon net.Conn
var log = logger.Logger

func init() {
	flag.Parse()
}

func main() {
	var err error
	waiCon, err = net.Dial("tcp", "127.0.0.1:8080")
	neiCon, err = net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Errorf("dial failed %s", err)
	}

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := waiCon.Read(buffer)
			if err != nil {
				log.Warnf("Failed to read nei tunnel connection %v", err)
			}

			_, err = neiCon.Write(buffer[:n])
			if err != nil {
				log.Warnf("Failed to write wai tunnel connection %v", err)
			}
		}
	}()

}
