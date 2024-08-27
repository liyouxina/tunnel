package protocal

import (
	"bufio"
	"errors"
	tcpIO "github.com/liyouxina/tunnel/pkg/io"
	"github.com/liyouxina/tunnel/pkg/logger"
	"io"
	"strconv"
	"strings"
)

var log = logger.Logger

type HTTPProtocol struct {
}

func (h HTTPProtocol) Do(task *Task, reader *bufio.Reader, writer *bufio.Writer) error {
	reqBody := task.genReqBody()
	if err := tcpIO.WriteAll(reqBody, writer); err != nil {
		return err
	}
	respBody, err := tcpIO.ReadAll(TAIL, reader)
	if err != nil && err != io.EOF {
		return err
	}
	if respBody == nil || *respBody == "" {
		return errors.New("respBody is nil")
	}
	res := strings.Split(*respBody, PARAM_SPLIT)
	task.ResStatus, _ = strconv.Atoi(res[0])
	if len(res) > 1 {
		task.ResBody = res[1]
	}
	return nil
}
