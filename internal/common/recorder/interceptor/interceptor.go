// Package interceptor 拦截器
package interceptor

import "github.com/ouqiang/mars/internal/common/recorder"

// Handlers 拦截器handlers
var Handlers []recorder.Interceptor

// 注册拦截器
func register(i recorder.Interceptor) {
	if i == nil {
		panic("interceptor is nil")
	}
	Handlers = append(Handlers, i)
}
