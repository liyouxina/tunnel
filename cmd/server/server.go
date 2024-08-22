package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/liyouxina/tunnel/pkg/protocal"
	"github.com/liyouxina/tunnel/pkg/tunnel"
	"net"
	"sync"

	"github.com/liyouxina/tunnel/pkg/logger"
)

var serverPort = flag.String("serverPort", "8090", "serverPort")
var tunnelPort = flag.String("tunnelPort", "8091", "tunnelPort")

var log = logger.Logger
var proto = protocal.HTTPProtocol{}
var taskPool = make(chan *protocal.Task, 100000)

func init() {
	flag.Parse()
}

func main() {
	startTunnelServer()
	startHTTPServer()
}

func startTunnelServer() {
	go func() {
		listener, err := net.Listen("tcp", ":"+*tunnelPort)
		if err != nil {
			log.Fatalf("start tunnel server failed %v", err)
			return
		}
		log.Infof("start tunnel server success %v", *tunnelPort)
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Warnf("Failed to accept tunnel connection %v", err)
				continue
			}
			tunnel.NewTunnel(conn, proto).Run(taskPool)
			log.Infof("get tunnel %s", conn.RemoteAddr())
		}
	}()
}

func startHTTPServer() {
	server := gin.Default()
	server.Any("/*path", func(c *gin.Context) {
		body := make([]byte, 8096)
		n, _ := c.Request.Body.Read(body)
		bodyString := string(body[:n])
		headersString := ""
		for k, v := range c.Request.Header {
			headersString = headersString + k + protocal.HEADER_SPLIT + v[0] + protocal.HEADER_TAIL
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		task := protocal.Task{
			Headers: headersString,
			Url:     c.Request.RequestURI,
			Body:    bodyString,
			Method:  c.Request.Method,
			Wg:      &wg,
		}
		taskPool <- &task
		wg.Wait()
		c.Status(task.ResStatus)
		_, _ = c.Writer.Write([]byte(task.ResBody))
	})
	_ = server.Run("0.0.0.0:" + *serverPort)
}
