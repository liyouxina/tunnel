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

GET / HTTP/1.1
Host: 127.0.0.1:8080
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate
Accept-Language: en-US,en;q=0.9
Connection: close
Upgrade-Insecure-Requests: 1
