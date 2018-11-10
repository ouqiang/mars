// Binary mars HTTP(S)代理
package main

import (
	"github.com/ouqiang/mars/cmd"
	"github.com/ouqiang/mars/internal/common/version"
)

var (
	// AppVersion 应用版本
	AppVersion string
	// BuildDate 构建日期
	BuildDate string
	// GitCommit 最后提交的git commit
	GitCommit string
)

func main() {
	version.Init(AppVersion, BuildDate, GitCommit)
	cmd.Execute()
}
