package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var server = flag.String("server", "localhost:8080", "server")
var localServer = flag.String("localServer", "localhost:8080", "localServer")
var clientCount = flag.Int("clientCount", 100, "clientCount")

var wg *sync.WaitGroup

func main() {
	flag.Parse()
	wg = &sync.WaitGroup{}
	for i := 0; i < *clientCount; i++ {
		wg.Add(1)
		go runClient(i + 1)
	}
	wg.Wait()
}

func runClient(number int) {
	// 1. 拨号方式建立与服务端连接
	conn, err := net.Dial("tcp", *server)
	if err != nil {
		fmt.Println("连接服务端失败,err:", err)
		wg.Done()
		return
	}

	// 注意：关闭连接位置，不能写在连接失败判断上面
	defer func(conn net.Conn) {
		wg.Done()
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)
	log.Println(`成功启动长连接` + strconv.Itoa(number))
	for {
		requestBody := make([]byte, 1024*8)
		n, _ := conn.Read(requestBody)
		requestBodyString := string(requestBody[:n])
		req := strings.Split(requestBodyString, `"""split"""`)
		headersString := req[0]
		headers := strings.Split(headersString, `;;;;;;;;;;;;;;;;;`)
		url := req[1]
		body := req[2]
		method := req[3]
		request, _ := http.NewRequest(method, "http://"+*localServer+url, bytes.NewBuffer([]byte(body)))
		for _, header := range headers {
			vs := strings.Split(header, `::::::::::::::::`)
			if len(vs) == 2 {
				request.Header.Set(vs[0], vs[1])
			}
		}
		resp, _ := http.DefaultClient.Do(request)
		respBody := make([]byte, 0, 1024*1024)
		respBody = append(respBody, strconv.Itoa(resp.StatusCode)...)
		respBody = append(respBody, []byte(`"""split"""`)...)
		respBodyString, _ := io.ReadAll(resp.Body)
		respBody = append(respBody, respBodyString...)
		_, _ = conn.Write(respBody)
	}
}
