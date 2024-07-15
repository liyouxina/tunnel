package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
)

type Tunnel struct {
	Ip   string
	conn net.Conn
}

func (tunnel *Tunnel) runTask() {
	reader := bufio.NewReader(tunnel.conn)
	for {
		// 读取客户端数据
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		// 输出客户端发送的数据
		log.Printf("Received from client: %s", message)

		// 向客户端发送响应
		response := fmt.Sprintf("Echo: %s", message)
		_, err = tunnel.conn.Write([]byte(response))
		if err != nil {
			log.Printf("Failed to write to connection: %v", err)
			return
		}
	}
}

func main() {
	startTunnels()
}

func startTunnels() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("Failed to close listener: %v", err)
		}
	}(listener)
	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		tunnel := &Tunnel{
			conn: conn,
			Ip:   conn.RemoteAddr().String(),
		}

		go tunnel.runTask()
	}
}

func startServer() {
	server :=
}
