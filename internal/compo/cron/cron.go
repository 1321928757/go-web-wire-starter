package cron

import (
	"context"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"go-web-wire-starter/internal/dao"
	"go.uber.org/zap"
)

// ProviderSet is cron providers.
var ProviderSet = wire.NewSet(NewCron, NewExampleJob)

type Cron struct {
	logger *zap.Logger
	data   *dao.Data
	server *cron.Cron

	exampleJob *ExampleJob
}

// NewCron .
func NewCron(data *dao.Data, logger *zap.Logger, exampleJob *ExampleJob) *Cron {
	server := cron.New(
		//让cron实例支持以秒为单位的定时任务。
		cron.WithSeconds(),
	)

	return &Cron{
		logger: logger,
		data:   data,
		server: server,

		exampleJob: exampleJob,
	}
}

func (c *Cron) Run() error {
	//cron example
	//if _, err := c.server.AddFunc("*/5 * * * * *", c.exampleJob.Hello); err != nil {
	//	return err
	//}

	c.server.Start()
	return nil
}

func (c *Cron) Stop(ctx context.Context) error {
	c.server.Stop()
	return nil
}
