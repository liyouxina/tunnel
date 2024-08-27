package io

import (
	"bufio"
	"strings"
)

func WriteAll(body []byte, writer *bufio.Writer) error {
	writeIndex := 0
	totalWriteLength := len(body)
	for writeIndex < totalWriteLength {
		n, err := writer.Write(body[writeIndex:])
		if err != nil {
			return err
		}
		writeIndex += n
	}
	_ = writer.Flush()
	return nil
}

func ReadAll(tail string, reader *bufio.Reader) (*string, error) {
	var respBodyString string
	respBody := make([]byte, 1024*1024*10)
	readIndex := 0
	for {
		n, err := reader.Read(respBody[readIndex:])
		if err != nil {
			return nil, err
		}
		readIndex += n
		respBodyString = string(respBody[:readIndex])
		tailIndex := strings.Index(respBodyString, tail)
		if tailIndex > 0 {
			respBodyString = respBodyString[:tailIndex]
			// 粘包处理 单tunnel串行处理请求，不会有粘包问题
			break
		}
	}
	return &respBodyString, nil
}
