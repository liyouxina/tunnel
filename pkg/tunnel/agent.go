package tunnel

import (
	"bufio"
	"bytes"
	tcpIO "github.com/liyouxina/tunnel/pkg/io"
	"github.com/liyouxina/tunnel/pkg/protocal"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type TunnelAgent struct {
	conn       net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	wg         *sync.WaitGroup
	httpServer string
}

func NewTunnelAgent(conn net.Conn, httpServer string, wg *sync.WaitGroup) *TunnelAgent {
	return &TunnelAgent{
		conn:       conn,
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
		httpServer: httpServer,
		wg:         wg,
	}
}

func (tunnelAgent *TunnelAgent) Run() {
	go func() {
		for {
			reqBody, err := tcpIO.ReadAll(protocal.TAIL, tunnelAgent.reader)
			if err != nil {
				log.Error("read reqBody fail", err)
				_ = tunnelAgent.conn.Close()
				tunnelAgent.wg.Done()
				break
			}
			if reqBody == nil || *reqBody == "" {
				log.Error("read reqBody nil or empty", err)
				_ = tunnelAgent.conn.Close()
				tunnelAgent.wg.Done()
				break
			}
			log.Infof("reqBody %s", *reqBody)
			reqParams := strings.Split(*reqBody, protocal.PARAM_SPLIT)

			headers := strings.Split(reqParams[0], protocal.HEADER_TAIL)
			url := reqParams[1]
			body := reqParams[2]
			method := reqParams[3]
			httpReq, _ := http.NewRequest(method, "http://"+tunnelAgent.httpServer+url, bytes.NewBuffer([]byte(body)))
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
				_ = tunnelAgent.conn.Close()
				tunnelAgent.wg.Done()
			}
		}
	}()
}
