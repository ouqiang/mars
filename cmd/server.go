package cmd

import (
	"path/filepath"
	"strings"

	"github.com/ouqiang/goutil"
	"github.com/ouqiang/mars/internal/app"
	"github.com/ouqiang/mars/internal/app/config"
	"github.com/ouqiang/mars/internal/app/inject"
	"github.com/ouqiang/mars/internal/common/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// 环境变量前缀
	serverEnvPrefix = "MARS"
	// 环境变量key分隔符
	serverConfigKeySeparator = "_"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run http server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info(version.Format())
		viper.BindPFlags(cmd.Flags())
		conf := createConfig()
		if conf.App.Env.IsDev() {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		container := inject.NewContainer(conf)
		app.New(container).Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	var env string
	var configFile string
	serverCmd.Flags().StringVarP(&env, "env", "e", "prod", "dev | prod")
	serverCmd.Flags().StringVarP(&configFile, "configFile", "c", "conf/app.toml", "config file path")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(serverEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", serverConfigKeySeparator))
	viper.SetConfigType("toml")
}

// 创建配置
func createConfig() *config.Config {
	currentDir, err := goutil.WorkDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := viper.GetString("configFile")
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(currentDir, configFile)
	}
	viper.SetConfigFile(configFile)
	log.Debugf("环境变量前缀: %s", serverEnvPrefix)
	log.Debugf("环境变量key分隔符: %s", serverConfigKeySeparator)
	log.Debugf("配置文件: %s", configFile)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("加载配置文件错误: %s", err)
	}
	conf := new(config.Config)
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("配置文件解析错误: %s", err)
	}
	conf.App.Env = config.RuntimeMode(viper.GetString("env"))

	return conf
}
