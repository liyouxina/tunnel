package protocal

import (
	"bufio"
	"sync"
)

const (
	HEADER_SPLIT = ":::"
	HEADER_TAIL  = ";;;"
	PARAM_SPLIT  = "<<<>>>"
	TAIL         = "=========================="
)

type Task struct {
	Headers   string
	Url       string
	Body      string
	Method    string
	Wg        *sync.WaitGroup
	Retry     int8
	ResStatus int
	ResBody   string
}

func (task *Task) genReqBody() []byte {
	reqBody := make([]byte, 0, 1024*1024)
	reqBody = append(reqBody, task.Headers...)
	reqBody = append(reqBody, []byte(PARAM_SPLIT)...)
	reqBody = append(reqBody, task.Url...)
	reqBody = append(reqBody, []byte(PARAM_SPLIT)...)
	reqBody = append(reqBody, task.Body...)
	reqBody = append(reqBody, []byte(PARAM_SPLIT)...)
	reqBody = append(reqBody, task.Method...)
	reqBody = append(reqBody, []byte(TAIL)...)
	return reqBody
}

type Protocol interface {
	Do(task *Task, reader *bufio.Reader, writer *bufio.Writer) error
}
