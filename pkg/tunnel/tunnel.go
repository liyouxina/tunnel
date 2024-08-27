package tunnel

import (
	"bufio"
	"github.com/liyouxina/tunnel/pkg/logger"
	"github.com/liyouxina/tunnel/pkg/protocal"
	"net"
)

var log = logger.Logger

type Tunnel struct {
	conn   net.Conn
	proto  protocal.Protocol
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewTunnel(conn net.Conn, proto protocal.Protocol) *Tunnel {
	return &Tunnel{
		conn:   conn,
		proto:  proto,
		reader: bufio.NewReaderSize(conn, 1024*1024*10),
		writer: bufio.NewWriter(conn),
	}
}

func (tunnel *Tunnel) Run(taskPool chan *protocal.Task) {
	go func() {
		for {
			task := <-taskPool
			err := tunnel.proto.Do(task, tunnel.reader, tunnel.writer)
			if err != nil {
				log.Error("do error, kill this tunnel ", err)
				_ = tunnel.conn.Close()
				taskPool <- task
				break
			}
			task.Wg.Done()
		}
	}()
}
