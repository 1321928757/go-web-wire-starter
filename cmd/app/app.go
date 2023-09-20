package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go-web-wire-starter/config"
	validator2 "go-web-wire-starter/util/validator"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"strings"
)

// App 应用程序
type App struct {
	// 应用程序的配置信息
	conf *config.Configuration
	// 应用程序的日志记录器
	logger *zap.Logger
	// 应用程序的 http 服务器
	httpSrv *http.Server
}

func newHttpServer(
	conf *config.Configuration,
	router *gin.Engine,
) *http.Server {
	// 初始化校验规则
	initValidator()
	return &http.Server{
		Addr:    ":" + conf.App.Port,
		Handler: router,
	}
}

func newApp(
	conf *config.Configuration,
	logger *zap.Logger,
	httpSrv *http.Server,
) *App {
	return &App{
		conf:    conf,
		logger:  logger,
		httpSrv: httpSrv,
	}
}

// 初始化校验规则
func initValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		_ = v.RegisterValidation("mobile", validator2.ValidateMobile)

		// 注册自定义 tag 函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// 'vn' tag - ValidatorMessages key name
			name := strings.SplitN(fld.Tag.Get("vn"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// 启动服务
func (a *App) Run() error {
	// 启动 http server
	go func() {
		a.logger.Info("http server started")
		if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	//启动其他服务，如cron定时器，queue队列消费者等

	return nil
}

// 关闭服务
func (a *App) Stop(ctx context.Context) error {
	// 关闭 http server
	a.logger.Info("http server has been stop")
	if err := a.httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	//关闭其他服务，如cron定时器，queue队列消费者等

	return nil
}
