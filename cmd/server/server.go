package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var serverPort = flag.String("serverPort", "8080", "serverPort")
var tunnelPort = flag.String("tunnelPort", "8080", "tunnelPort")

type Task struct {
	headers   map[string]*string
	url       string
	body      string
	method    string
	wg        *sync.WaitGroup
	resStatus int
	resBody   string
}

var taskPool chan *Task

type Tunnel struct {
	Ip   string
	conn net.Conn
}

func (tunnel *Tunnel) runTask() {
	for {
		task := <-taskPool
		conn := tunnel.conn
		taskBody := make([]byte, 8096)
		headerBody, _ := json.Marshal(task.headers)
		taskBody = append(taskBody, headerBody...)
		taskBody = append(taskBody, []byte(`"""split"""`)...)
		taskBody = append(taskBody, task.url...)
		taskBody = append(taskBody, []byte(`"""split"""`)...)
		taskBody = append(taskBody, task.body...)
		taskBody = append(taskBody, []byte(`"""split"""`)...)
		taskBody = append(taskBody, task.method...)
		_, _ = conn.Write(taskBody)
		reader := bufio.NewReader(tunnel.conn)
		respBody := make([]byte, 1024*1024)
		n, _ := reader.Read(respBody)
		respBodyString := string(respBody[:n])
		res := strings.Split(respBodyString, `"""split"""`)
		task.resStatus, _ = strconv.Atoi(res[0])
		if len(res) > 1 {
			task.resBody = res[1]
		}
		task.wg.Done()
	}

}

func startTunnels() {
	listener, err := net.Listen("tcp", ":"+*tunnelPort)
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
	server := gin.Default()
	server.Any("/*path", func(c *gin.Context) {
		body := make([]byte, 8096)
		n, _ := c.Request.Body.Read(body)
		bodyString := string(body[:n])
		wg := sync.WaitGroup{}
		wg.Add(1)
		task := Task{
			headers: make(map[string]*string),
			url:     c.Request.RequestURI,
			body:    bodyString,
			method:  c.Request.Method,
			wg:      &wg,
		}
		taskPool <- &task
		wg.Wait()
		c.Status(task.resStatus)
		_, _ = c.Writer.Write([]byte(task.resBody))
	})
	_ = server.Run("0.0.0.0:" + *serverPort)
}

func main() {
	taskPool = make(chan *Task)
	flag.Parse()
	go startTunnels()
	startServer()
}
