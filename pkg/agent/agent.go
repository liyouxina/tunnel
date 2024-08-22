package agent

import (
	"bufio"
	"bytes"
	tcpIO "github.com/liyouxina/tunnel/pkg/io"
	"github.com/liyouxina/tunnel/pkg/logger"
	"github.com/liyouxina/tunnel/pkg/protocal"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var log = logger.Logger

var tunnelMakeSignal chan int

func StartAgents(tunnelServer string, targetServer string, maxAgentCnt int) {
	tunnelMakeSignal = make(chan int, maxAgentCnt)
	for i := 0; i < maxAgentCnt; i++ {
		tunnelMakeSignal <- 1
	}

	go func() {
		for {
		_:
			<-tunnelMakeSignal

			conn, err := net.Dial("tcp", tunnelServer)
			if err != nil {
				log.Errorf("dial failed %s", err)
				tunnelMakeSignal <- 1
			}

			tunnelAgent := &TunnelAgent{
				conn:         conn,
				reader:       bufio.NewReader(conn),
				writer:       bufio.NewWriter(conn),
				targetServer: targetServer,
			}
			tunnelAgent.Run()
			log.Infof(`open tunnel agent`)
		}
	}()
}

type TunnelAgent struct {
	conn         net.Conn
	reader       *bufio.Reader
	writer       *bufio.Writer
	targetServer string
}

func (tunnelAgent *TunnelAgent) Run() {
	go func() {
		defer func() {
			tunnelAgent.Close()
		}()
		for {
			reqBody, err := tcpIO.ReadAll(protocal.TAIL, tunnelAgent.reader)
			if err != nil || reqBody == nil || *reqBody == "" {
				log.Error("read reqBody fail", err, reqBody)
				tunnelAgent.Close()
				break
			}
			log.Infof("reqBody %s", *reqBody)
			reqParams := strings.Split(*reqBody, protocal.PARAM_SPLIT)
			url := reqParams[1]
			body := reqParams[2]
			method := reqParams[3]
			httpReq, _ := http.NewRequest(method, "http://"+tunnelAgent.targetServer+url, bytes.NewBuffer([]byte(body)))
			headers := strings.Split(reqParams[0], protocal.HEADER_TAIL)
			for _, header := range headers {
				vs := strings.Split(header, protocal.HEADER_SPLIT)
				if len(vs) == 2 {
					httpReq.Header.Set(vs[0], vs[1])
				}
			}
			httpResp, _ := http.DefaultClient.Do(httpReq)
			respBody := make([]byte, 0, 1024*1024)
			respBody = append(respBody, strconv.Itoa(httpResp.StatusCode)...)
			respBody = append(respBody, []byte(protocal.PARAM_SPLIT)...)
			respBodyString, _ := io.ReadAll(httpResp.Body)
			respBody = append(respBody, respBodyString...)
			respBody = append(respBody, []byte(protocal.TAIL)...)
			log.Infof("respBody %s", respBody)
			if err = tcpIO.WriteAll(respBody, tunnelAgent.writer); err != nil {
				log.Error("write respBody fail", err)
				tunnelAgent.Close()
			}
		}
	}()
}

func (tunnelAgent *TunnelAgent) Close() {
	log.Info("close tunnel agent")
	_ = tunnelAgent.conn.Close()
	tunnelMakeSignal <- 1
}
