// package inject 依赖注入
package inject

import (
	"os"
	"path/filepath"

	"github.com/ouqiang/mars/internal/common/recorder/interceptor"

	"github.com/ouqiang/mars/internal/common/recorder/output"

	"github.com/ouqiang/goproxy"
	"github.com/ouqiang/mars/internal/app/config"
	"github.com/ouqiang/mars/internal/common"
	"github.com/ouqiang/mars/internal/common/recorder"
	"github.com/ouqiang/mars/internal/common/recorder/storage"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

// Container 容器
type Container struct {
	Conf           *config.Config
	Proxy          *goproxy.Proxy
	TxStorage      recorder.Storage
	txRecorder     *recorder.Recorder
	txOutput       recorder.Output
	txInterceptors []recorder.Interceptor
}

// NewContainer 创建容器
func NewContainer(conf *config.Config) *Container {
	if conf == nil {
		panic("config is nil")
	}
	c := &Container{
		Conf:       conf,
		txRecorder: recorder.NewRecorder(),
	}
	c.createProxy()
	c.createRecorderStorage()
	c.createRecorderOutput()
	c.createRecorderInterceptor()

	c.txRecorder.SetProxy(c.Proxy)
	c.txRecorder.SetStorage(c.TxStorage)
	c.txRecorder.SetOutput(c.txOutput)
	c.txRecorder.SetInterceptors(c.txInterceptors)

	return c
}

func (c *Container) createProxy() {
	opts := make([]goproxy.Option, 0, 3)
	opts = append(opts, goproxy.WithDisableKeepAlive(true))
	if c.Conf.MITMProxy.Enabled {
		opts = append(opts, goproxy.WithDelegate(c.txRecorder))
	}
	if c.Conf.MITMProxy.DecryptHTTPS {
		queue := common.NewQueue(c.Conf.MITMProxy.CertCacheSize)
		certCache := recorder.NewCertCache(queue)
		opts = append(opts, goproxy.WithDecryptHTTPS(certCache))
	}
	c.Proxy = goproxy.New(opts...)
}

func (c *Container) createRecorderStorage() {
	if !c.Conf.MITMProxy.Enabled {
		return
	}
	if c.Conf.MITMProxy.LeveldbDir == "" {
		c.Conf.MITMProxy.LeveldbDir = filepath.Join(os.TempDir(), "mars_leveldb")
	}
	if _, err := os.Stat(c.Conf.MITMProxy.LeveldbDir); err == nil {
		err = os.RemoveAll(c.Conf.MITMProxy.LeveldbDir)
		if err != nil {
			log.Fatalf("删除leveldb数据库目录错误: %s", err)
		}
	}

	db, err := leveldb.OpenFile(c.Conf.MITMProxy.LeveldbDir, nil)
	if err != nil {
		log.Fatalf("创建leveldb数据库错误: %s", err)
	}
	queue := common.NewQueue(c.Conf.MITMProxy.LeveldbCacheSize)
	c.TxStorage = storage.NewLevelDB(db, queue)
}

func (c *Container) createRecorderOutput() {
	c.txOutput = output.NewConsole()
}

func (c *Container) createRecorderInterceptor() {
	c.txInterceptors = interceptor.Handlers
}
