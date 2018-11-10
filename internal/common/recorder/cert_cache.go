package recorder

import (
	"crypto/tls"
	"sync"

	"github.com/ouqiang/mars/internal/common"
)

// CertCache 证书缓存
type CertCache struct {
	queue *common.Queue
	m     sync.Map
}

// 创建CertCache
func NewCertCache(q *common.Queue) *CertCache {
	c := &CertCache{
		queue: q,
	}

	return c
}

// Get 获取证书
func (c *CertCache) Get(host string) *tls.Certificate {
	value, ok := c.m.Load(host)
	if !ok {
		return nil
	}

	return value.(*tls.Certificate)
}

// Set 保存证书
func (c *CertCache) Set(host string, cert *tls.Certificate) {
	removeValue := c.queue.Add(host)
	if removeValue != nil {
		c.m.Delete(removeValue.(string))
	}
	c.m.Store(host, cert)
}
