package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-web-wire-starter/config"
	"go-web-wire-starter/internal/command"
	"go-web-wire-starter/util/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	// rootPath 项目根目录
	rootPath = path.RootPath()
	// Version 版本号
	Version string
	// configPath 配置文件路径
	configPath string
	// conf 配置文件
	conf         *config.Configuration
	loggerWriter *lumberjack.Logger
	// logger 日志
	logger *zap.Logger
)

func init() {
	//  解析命令行参数conf(配置文件，满足生产环境与开发环境的需求)，如果没有指定则使用默认值config.yaml
	pflag.StringVarP(&configPath, "conf", "", filepath.Join(rootPath, "conf", "config.yaml"), "config path, eg: --conf config.yaml")
	//  初始化日志和配置文件
	cobra.OnInitialize(func() {
		initConfig()
		initLogger()
	})
}

func main() {
	// 创建的命令行程序的根命令
	rootCmd := &cobra.Command{
		// 命令的名称
		Use: "internal",
		// 执行命令时会调用此函数
		Run: func(cmd *cobra.Command, args []string) {
			app, cleanup, err := wireApp(conf, loggerWriter, logger)
			if err != nil {
				panic(err)
			}
			defer cleanup()

			// 启动应用
			log.Printf("start internal %s ...", Version)
			if err := app.Run(); err != nil {
				panic(err)
			}

			// 等待中断信号以优雅地关闭应用
			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			log.Printf("shutdown internal %s ...", Version)

			// 设置 5 秒的超时时间
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// 关闭应用
			if err := app.Stop(ctx); err != nil {
				panic(err)
			}
		},
	}

	// 将待注册的子命令注册到根命令
	command.Register(rootCmd, func() (*command.Command, func(), error) {
		return wireCommand(conf, loggerWriter, logger)
	})

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// initConfig 初始化配置文件
func initConfig() {
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(rootPath, "conf", configPath)
	}

	fmt.Println("load config:" + configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file_and_disk changed:", in.Name)
		defer func() {
			if err := recover(); err != nil {
				logger.Error("config file_and_disk changed err:", zap.Any("err", err))
				fmt.Println(err)
			}
		}()
		if err := v.Unmarshal(&conf); err != nil {
			fmt.Println(err)
		}
	})
}

// initLogger 初始化日志
func initLogger() {
	var level zapcore.Level  // zap 日志等级
	var options []zap.Option // zap 配置项

	logFileDir := conf.Log.RootDir
	if !filepath.IsAbs(logFileDir) {
		logFileDir = filepath.Join(rootPath, logFileDir)
	}

	if ok, _ := path.Exists(logFileDir); !ok {
		_ = os.Mkdir(conf.Log.RootDir, os.ModePerm)
	}

	switch conf.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(conf.App.Env + "." + l.String())
	}

	loggerWriter = &lumberjack.Logger{
		Filename:   filepath.Join(logFileDir, conf.Log.Filename),
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   conf.Log.Compress,
	}

	logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(loggerWriter), level), options...)
}
