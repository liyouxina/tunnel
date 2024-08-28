package main

import (
	"flag"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"

	"github.com/liyouxina/tunnel/pkg/logger"
)

var serverPort = flag.String("serverPort", "8089", "serverPort")
var tunnelPort = flag.String("tunnelPort", "8091", "tunnelPort")
var serverConn net.Conn
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
	listener, err := net.Listen("tcp", ":"+*tunnelPort)
	if err != nil {
		log.Error("start tunnel server failed", zap.Error(err))
		return
	}
	log.Info("start tunnel server success", zap.String("port", *tunnelPort))
	go func() {
		for {
			tunnelConn, err = listener.Accept()
			if err != nil {
				log.Error("Failed to accept tunnel connection", zap.Error(err))
				continue
			}
			log.Info("get tunnelConn from", zap.String("client", tunnelConn.RemoteAddr().String()))

			go func() {
				for {
					buffer := make([]byte, 1024)
					for tunnelConn == nil {
						log.Info("waiting for tunnel connection")
						time.Sleep(time.Second)
					}
					log.Info("waiting for tunnel connection read")
					n, err := tunnelConn.Read(buffer)
					log.Info("tunnel connection read success", zap.String("client", string(buffer[:n])))
					if err != nil {
						log.Error("Failed to read tunnelConn", zap.Error(err))
					}
					for serverConn == nil {
						log.Info("waiting for serverConn from", zap.String("client", tunnelConn.RemoteAddr().String()))
						time.Sleep(time.Second)
					}
					if n != 0 {
						log.Info("write to serverConn", zap.String("client", string(buffer[:n])))
						_, err = serverConn.Write(buffer[:n])
						if err != nil {
							log.Error("Failed to write serverConn", zap.Error(err))
						}
					}
					time.Sleep(time.Second)
				}

			}()
		}
	}()

}

var lock sync.Mutex

func server() {
	var err error
	listener, err := net.Listen("tcp", ":"+*serverPort)
	if err != nil {
		log.Error("start server failed", zap.Error(err))
		return
	}
	log.Info("start server success")
	go func() {
		for {
			serverConn, err = listener.Accept()
			if err != nil {
				log.Error("Failed to accept server connection", zap.Error(err))
				continue
			}
			log.Info("get serverConn from client", zap.String("client", serverConn.RemoteAddr().String()))
			lock.Lock()
			go func() {
				for {
					buffer := make([]byte, 1024)
					for serverConn == nil {
						log.Info("waiting for serverConn")
						time.Sleep(time.Second)
					}
					log.Info("waiting for serverConn read")
					n, err := serverConn.Read(buffer)
					log.Info("serverConn read success", zap.String("client", string(buffer[:n])))
					if err != nil {
						log.Error("Failed to read serverConn", zap.Error(err))
						err = serverConn.Close()
						lock.Unlock()
						if err != nil {
							log.Error("Failed to close server connection", zap.Error(err))
						} else {
							log.Info("close server connection success")
						}
						break
					}
					for tunnelConn == nil {
						log.Info("waiting for tunnelConn")
						time.Sleep(time.Second)
					}
					if n != 0 {
						log.Info("write to tunnelConn", zap.String("client", string(buffer[:n])))
						_, err = tunnelConn.Write(buffer[:n])
						if err != nil {
							log.Error("Failed to write tunnelConn", zap.Error(err))
							_ = tunnelConn.Close()
							break
						}
					}
					time.Sleep(time.Second)
				}

			}()
		}
	}()
}
