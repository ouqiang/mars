// Package config 配置
package config

import (
	"net"
	"strconv"
)

// RuntimeMode 运行模式
type RuntimeMode string

func (m RuntimeMode) IsDev() bool {
	return m == "dev"
}

func (m RuntimeMode) IsProd() bool {
	return m == "prod"
}

// Config 配置
type Config struct {
	// App 应用配置
	App appConfig `mapstructure:"app"`
	// Proxy 代理配置
	MITMProxy mitmProxyConfig `mapstructure:"mitmProxy"`
}

type appConfig struct {
	Env           RuntimeMode
	Host          string `mapstructure:"host"`
	ProxyPort     int    `mapstructure:"proxyPort"`
	InspectorPort int    `mapstructure:"inspectorPort"`
}

type mitmProxyConfig struct {
	Enabled          bool   `mapstructure:"enabled"`
	DecryptHTTPS     bool   `mapstructure:"decryptHTTPS"`
	CertCacheSize    int    `mapstructure:"certCacheSize"`
	LeveldbDir       string `mapstructure:"leveldbDir"`
	LeveldbCacheSize int    `mapstructure:"leveldbCacheSize"`
}

// ProxyAddr 代理监听地址
func (ac appConfig) ProxyAddr() string {
	return net.JoinHostPort(ac.Host, strconv.Itoa(ac.ProxyPort))
}

// InspectorAddr 审查监听地址
func (ac appConfig) InspectorAddr() string {
	return net.JoinHostPort(ac.Host, strconv.Itoa(ac.InspectorPort))
}
