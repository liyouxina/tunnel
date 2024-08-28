package main

import (
	"flag"
	"github.com/liyouxina/tunnel/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
)

var targetServer = flag.String("targetServer", "127.0.0.1:8066", "targetServer")
var tunnelServer = flag.String("tunnelServer", "127.0.0.1:8091", "tunnelServer")
var targetConn net.Conn
var tunnelConn net.Conn
var log = logger.GetLogger()

func init() {
	flag.Parse()
}

func main() {
	tunnel()
	server()
	for {
		time.Sleep(time.Minute)
	}
}

func tunnel() {
	var err error
	tunnelConn, err = net.Dial("tcp", *tunnelServer)
	if err != nil {
		log.Error("connet tunnel server failed", zap.Error(err))
		return
	}
	log.Info("connet tunnel server success", zap.String("server", *tunnelServer))
	go func() {
		buffer := make([]byte, 1024)
		for {
			log.Info("wait for tunnelConn read")
			n, err := tunnelConn.Read(buffer)
			log.Info("tunnelConn read success", zap.String("content", string(buffer[:n])))
			if err != nil && err != io.EOF {
				log.Error("Failed to read tunnelConn", zap.Error(err))
				_ = tunnelConn.Close()
				break
			}

			if targetConn == nil {
				targetConn, err = net.Dial("tcp", *targetServer)
				if err != nil {
					log.Error("Failed to dial targetServer", zap.Error(err))
					break
				}
			}
			if n == 0 {
				time.Sleep(time.Second)
				continue
			}
			log.Info("targetConn write", zap.String("content", string(buffer[:n])))
			_, err = targetConn.Write(buffer[:n])
			if err != nil {
				log.Error("Failed to write tunnelConn", zap.Error(err))
				_ = targetConn.Close()
				targetConn = nil
			}
			time.Sleep(time.Second)
		}
	}()

}

func server() {
	go func() {
		buffer := make([]byte, 1024)
		for {
			for targetConn == nil {
				time.Sleep(time.Second)
			}
			log.Info("wait for targetConn read")
			n, err := targetConn.Read(buffer)
			log.Info("targetConn read success", zap.String("content", string(buffer[:n])))
			if err != nil {
				log.Error("Failed to read tunnelConn", zap.Error(err))
				if err != io.EOF {
					_ = targetConn.Close()
					targetConn = nil
				}
			}
			for tunnelConn == nil {
				log.Info("wait for tunnelConn")
				time.Sleep(time.Second)
			}
			if n == 0 {
				time.Sleep(time.Second)
				continue
			}
			log.Info("tunnelConn write", zap.String("content", string(buffer[:n])))
			_, err = tunnelConn.Write(buffer[:n])
			if err != nil {
				log.Error("Failed to write tunnelConn", zap.Error(err))
				_ = tunnelConn.Close()
			}
			time.Sleep(time.Second)
		}

	}()

}
