package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端请求
	reader := bufio.NewReader(conn)
	request, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	fmt.Println("Received request:", request)

	// 解析请求行
	requestLine := strings.Fields(request)
	if len(requestLine) < 3 {
		fmt.Println("Malformed request")
		return
	}
	method := requestLine[0]
	uri := requestLine[1]

	fmt.Printf("Method: %s, URI: %s\n", method, uri)

	// 构建 HTTP 响应
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		"<html><body><h1>Hello, World!</h1></body></html>"

	// 发送响应给客户端
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}

func main() {
	// 监听 TCP 连接
	listener, err := net.Listen("tcp", ":8034")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 处理客户端连接
		go handleConnection(conn)
	}
}
