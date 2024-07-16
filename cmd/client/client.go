package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/liyouxina/tunnel/pkg/logger"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
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
	log.Infof(`tcp tunnel start ` + strconv.Itoa(number))
	for {
		requestBody := make([]byte, 1024*8)
		n, _ := conn.Read(requestBody)
		requestBodyString := string(requestBody[:n])
		log.Infof("req body %s", requestBodyString)
		reqParams := strings.Split(requestBodyString, `"""split"""`)
		headersString := reqParams[0]
		headers := strings.Split(headersString, `;;;;;;;;;;;;;;;;;`)
		url := reqParams[1]
		body := reqParams[2]
		method := reqParams[3]
		httpRequest, _ := http.NewRequest(method, "http://"+*localServer+url, bytes.NewBuffer([]byte(body)))
		for _, header := range headers {
			vs := strings.Split(header, `::::::::::::::::`)
			if len(vs) == 2 {
				httpRequest.Header.Set(vs[0], vs[1])
			}
		}
		resp, _ := http.DefaultClient.Do(httpRequest)
		respBody := make([]byte, 0, 1024*1024)
		respBody = append(respBody, strconv.Itoa(resp.StatusCode)...)
		respBody = append(respBody, []byte(`"""split"""`)...)
		respBodyString, _ := io.ReadAll(resp.Body)
		respBody = append(respBody, respBodyString...)
		log.Infof("resp body %s", respBody)
		_, _ = conn.Write(respBody)
	}
}
