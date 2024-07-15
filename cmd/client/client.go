package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

/* TCP 客户端配置 */

func main() {
	// 1. 拨号方式建立与服务端连接
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("连接服务端失败,err:", err)
		return
	}

	// 注意：关闭连接位置，不能写在连接失败判断上面
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	for {
		requestBody := make([]byte, 1024*8)
		n, _ := conn.Read(requestBody)
		requestBodyString := string(requestBody[:n])
		req := strings.Split(requestBodyString, `"""split"""`)
		headersString := req[0]
		headers := make(map[string]*string)
		_ = json.Unmarshal([]byte(headersString), &headers)
		url := req[1]
		body := req[2]
		method := req[3]
		request, _ := http.NewRequest(method, "http://localhost:8082/"+url, bytes.NewBuffer([]byte(body)))
		for k, v := range headers {
			request.Header.Set(k, *v)
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
