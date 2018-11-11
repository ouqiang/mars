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
func (c *Console) Write(tx *recorder.Transaction) error {
	c.m.Lock()
	defer c.m.Unlock()

	c.builder.WriteString("txId: ")
	c.builder.WriteString(tx.Id)
	c.builder.WriteString("\n")
	c.writeRequest(tx)
	c.builder.WriteString("\n")
	c.writeResponse(tx)
	c.builder.WriteString("\n\n\n")
	io.WriteString(c.writer, c.builder.String())
	c.builder.Reset()

	return nil
}

func (c *Console) writeRequest(tx *recorder.Transaction) {
	c.builder.WriteString(fmt.Sprintf("\033[34m %s \033[0m", tx.Req.Method))
	c.builder.WriteString("  ")
	c.builder.WriteString(tx.Req.URL)
	c.builder.WriteString("\n")
	c.builder.WriteString("Server-IP: ")
	c.builder.WriteString(tx.ServerIP)
	c.builder.WriteString("\n")
	for key, values := range tx.Req.Header {
		c.builder.WriteString(key)
		c.builder.WriteString(": ")
		c.builder.WriteString(strings.Join(values, ";"))
		c.builder.WriteString("\n")
	}
	if !tx.Req.Body.IsBinary {
		c.builder.Write(tx.Req.Body.Content)
		c.builder.WriteString("\n")
	}
}

func (c *Console) writeResponse(tx *recorder.Transaction) {
	if tx.Resp.Err != nil {
		c.builder.WriteString(tx.Resp.Err.Error())
		return
	}
	c.builder.WriteString(fmt.Sprintf("\033[32m %s %s \033[0m", tx.Resp.Proto, tx.Resp.Status))
	c.builder.WriteString("\n")
	for key, values := range tx.Resp.Header {
		c.builder.WriteString(key)
		c.builder.WriteString(": ")
		c.builder.WriteString(strings.Join(values, ";"))
		c.builder.WriteString("\n")
	}
	if !tx.Resp.Body.IsBinary {
		c.builder.Write(tx.Resp.Body.Content)
	}
}
