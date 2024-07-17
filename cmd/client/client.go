package main

import (
	"flag"
	"github.com/liyouxina/tunnel/pkg/logger"
	"github.com/liyouxina/tunnel/pkg/tunnel"
	"net"
	"strconv"
	"sync"
)

var server = flag.String("server", "localhost:8080", "server")
var localServer = flag.String("localServer", "localhost:8080", "localServer")
var clientCount = flag.Int("clientCount", 100, "clientCount")
var log = logger.Logger

var wg *sync.WaitGroup

func main() {
	flag.Parse()
	wg = &sync.WaitGroup{}
	for number := 0; number < *clientCount; number++ {
		wg.Add(1)
		conn, err := net.Dial("tcp", *server)
		if err != nil {
			log.Errorf("连接服务端失败 %s", err)
			wg.Done()
		}
		tunnelAgent := tunnel.NewTunnelAgent(conn, *localServer, wg)
		tunnelAgent.Run()
		log.Infof(`tcp tunnel start ` + strconv.Itoa(number))
	}
	wg.Wait()
}
