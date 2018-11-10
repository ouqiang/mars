// Package storage 存储http transaction
package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ouqiang/mars/internal/common/recorder"
)

// Console 输出到终端
type Console struct {
	m       sync.Mutex
	builder strings.Builder
	writer  io.Writer
}

// NewConsole 创建console
func NewConsole() *Console {
	return &Console{
		writer: os.Stdout,
	}
}

// Write transaction输出到终端
func (s *Console) Write(tx *recorder.Transaction) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.builder.WriteString("txId: ")
	s.builder.WriteString(tx.Id)
	s.builder.WriteString("\n")
	s.writeRequest(tx)
	s.builder.WriteString("\n")
	s.writeResponse(tx)
	s.builder.WriteString("\n\n\n")
	io.WriteString(s.writer, s.builder.String())
	s.builder.Reset()

	return nil
}

func (s *Console) writeRequest(tx *recorder.Transaction) {
	s.builder.WriteString(fmt.Sprintf("\033[34m %s \033[0m", tx.Req.Method))
	s.builder.WriteString("  ")
	s.builder.WriteString(tx.Req.URL)
	s.builder.WriteString("\n")
	s.builder.WriteString("Server-IP: ")
	s.builder.WriteString(tx.ServerIP)
	s.builder.WriteString("\n")
	for key, values := range tx.Req.Header {
		s.builder.WriteString(key)
		s.builder.WriteString(": ")
		s.builder.WriteString(strings.Join(values, ";"))
		s.builder.WriteString("\n")
	}
	if !tx.Req.Body.IsBinary {
		s.builder.Write(tx.Req.Body.Content)
		s.builder.WriteString("\n")
	}
}

func (s *Console) writeResponse(tx *recorder.Transaction) {
	if tx.Resp.Err != nil {
		s.builder.WriteString(tx.Resp.Err.Error())
		return
	}
	s.builder.WriteString(fmt.Sprintf("\033[32m %s %s \033[0m", tx.Resp.Proto, tx.Resp.Status))
	s.builder.WriteString("\n")
	for key, values := range tx.Resp.Header {
		s.builder.WriteString(key)
		s.builder.WriteString(": ")
		s.builder.WriteString(strings.Join(values, ";"))
		s.builder.WriteString("\n")
	}
	if !tx.Resp.Body.IsBinary {
		s.builder.Write(tx.Resp.Body.Content)
	}
}
