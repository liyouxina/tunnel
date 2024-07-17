package protocal

import (
	"bufio"
	"github.com/liyouxina/tunnel/pkg/io"
	"github.com/liyouxina/tunnel/pkg/logger"
	"strconv"
	"strings"
)

var log = logger.Logger

type HTTPProtocol struct {
}

func (h HTTPProtocol) Do(task *Task, reader *bufio.Reader, writer *bufio.Writer) error {
	reqBody := task.genReqBody()
	if err := io.WriteAll(reqBody, writer); err != nil {
		return err
	}
	log.Infof("reqBody %s", string(reqBody))
	respBody, err := io.ReadAll(TAIL, reader)
	if err != nil {
		return err
	}
	res := strings.Split(*respBody, PARAM_SPLIT)
	task.ResStatus, _ = strconv.Atoi(res[0])
	if len(res) > 1 {
		task.ResBody = res[1]
	}
	return nil
}
