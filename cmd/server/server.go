package main

import (
	"bufio"
	"flag"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/liyouxina/tunnel/pkg/logger"
)

var serverPort = flag.String("serverPort", "8080", "serverPort")
var tunnelPort = flag.String("tunnelPort", "8080", "tunnelPort")

var log = logger.Logger

type Task struct {
	headers   string
	url       string
	body      string
	method    string
	wg        *sync.WaitGroup
	resStatus int
	resBody   string
}

var taskPool = make(chan *Task, 100000)

type Tunnel struct {
	Ip     string
	Number int
	conn   net.Conn
}

func (tunnel *Tunnel) runTask() {
	for {
		task := <-taskPool
		conn := tunnel.conn
		taskBody := make([]byte, 0, 8096*10)
		taskBody = append(taskBody, task.headers...)
		taskBody = append(taskBody, []byte(`"""split"""`)...)
		taskBody = append(taskBody, task.url...)
		taskBody = append(taskBody, []byte(`"""split"""`)...)
		taskBody = append(taskBody, task.body...)
		taskBody = append(taskBody, []byte(`"""split"""`)...)
		taskBody = append(taskBody, task.method...)
		_, err := conn.Write(taskBody)
		if err != nil {
			log.Warnf("write error %v", err)
			log.Warnf(`close conn %s %d`, tunnel.conn.RemoteAddr(), tunnel.Number)
			taskPool <- task
			break
		}
		log.Infof("send taskBody %s", string(taskBody))
		reader := bufio.NewReader(tunnel.conn)
		respBody := make([]byte, 1024*1024)
		n, err := reader.Read(respBody)
		if err != nil {
			log.Warnf("read error %v", err)
			log.Warnf(`close conn %s %d`, tunnel.conn.RemoteAddr(), tunnel.Number)
			taskPool <- task
			break
		}
		respBodyString := string(respBody[:n])
		log.Infof("read taskResp %s", respBodyString)
		res := strings.Split(respBodyString, `"""split"""`)
		task.resStatus, _ = strconv.Atoi(res[0])
		if len(res) > 1 {
			task.resBody = res[1]
		}
		task.wg.Done()
	}

}

func startTunnelServer() {
	listener, err := net.Listen("tcp", ":"+*tunnelPort)
	if err != nil {
		log.Fatalf("start tunnel server failed %v", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Failed to close listener %v", err)
		}
	}(listener)
	log.Infof("start tunnel server success %v", *tunnelPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warnf("Failed to accept tunnel connection %v", err)
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
		headersString := ""
		for k, v := range c.Request.Header {
			headersString = headersString + k + "::::::::::::::::" + v[0] + ";;;;;;;;;;;;;;;;;"
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		task := Task{
			headers: headersString,
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
	flag.Parse()
	go startTunnelServer()
	startServer()
}
