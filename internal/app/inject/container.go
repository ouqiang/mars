// package inject 依赖注入
package inject

import "github.com/ouqiang/mars/internal/app/config"

// Container 容器
type Container struct {
	Conf *config.Config
}

// NewContainer 创建容器
func NewContainer(conf *config.Config) *Container {
	c := &Container{
		Conf: conf,
	}

	return c
}
