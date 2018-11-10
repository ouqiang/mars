// Package version 版本
package version

import "github.com/ouqiang/goutil"

var (
	appVersion string
	buildDate  string
	gitCommit  string
)

// Init 初始化版本信息
func Init(version, date, commitId string) {
	appVersion = version
	buildDate = date
	gitCommit = commitId
}

// Format 格式化版本信息
func Format() string {
	out, err := goutil.FormatAppVersion(appVersion, gitCommit, buildDate)
	if err != nil {
		panic(err)
	}

	return out
}
