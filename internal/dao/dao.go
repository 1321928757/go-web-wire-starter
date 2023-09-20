package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/sony/sonyflake"
	"go-web-wire-starter/config"
	"go-web-wire-starter/util/path"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// ProviderSet Provider对象集合
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewUserDao, NewJwtDao)

type Data struct {
	db  *gorm.DB
	rdb *redis.Client
	sf  *sonyflake.Sonyflake
}

// NewData .
func NewData(logger *zap.Logger, db *gorm.DB, rdb *redis.Client, sf *sonyflake.Sonyflake) (*Data, func(), error) {
	cleanup := func() {
		logger.Info("closing the data resources")
	}

	return &Data{db: db, rdb: rdb, sf: sf}, cleanup, nil
}

// NewDB .
func NewDB(conf *config.Configuration, gLog *zap.Logger) *gorm.DB {
	if conf.Database.Driver != "mysql" {
		panic(conf.Database.Driver + " driver is not supported")
	}

	var writer io.Writer
	var logMode logger.LogLevel

	// 是否启用日志文件
	if conf.Database.EnableFileLogWriter {
		logFileDir := conf.Log.RootDir
		if !filepath.IsAbs(logFileDir) {
			logFileDir = filepath.Join(path.RootPath(), logFileDir)
		}
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   filepath.Join(logFileDir, conf.Database.LogFilename),
			MaxSize:    conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			MaxAge:     conf.Log.MaxAge,
			Compress:   conf.Log.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}

	switch conf.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                        // 慢查询 SQL 阈值
			Colorful:                  !conf.Database.EnableFileLogWriter, // 禁用彩色打印
			IgnoreRecordNotFoundError: false,                              // 忽略ErrRecordNotFound（记录未找到）错误
			LogLevel:                  logMode,                            // Log lever
		},
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Database.UserName,
		conf.Database.Password,
		conf.Database.Host,
		strconv.Itoa(conf.Database.Port),
		conf.Database.Database,
		conf.Database.Charset,
	)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Database.TablePrefix,
			//SingularTable: true,

		},
		DisableForeignKeyConstraintWhenMigrating: true,      // 禁用自动创建外键约束
		Logger:                                   newLogger, // 使用自定义 Logger
	}); err != nil {
		gLog.Error("failed opening connection to err:", zap.Any("err", err))
		panic("failed to connect database")
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
		return db
	}
}

// NewRedis .
func NewRedis(c *config.Configuration, gLog *zap.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + c.Redis.Port,
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	})

	//为 Redis 客户端添加一个 OpenTelemetry 的追踪钩子
	//OpenTelemetry 是用于分布式追踪的开源工具，用于监控和分析分布式系统中的请求流程
	client.AddHook(redisotel.TracingHook{})
	if err := client.Ping(context.Background()).Err(); err != nil {
		gLog.Error("redis connect failed, err:", zap.Any("err", err))
		panic("failed to connect redis")
	} else {
		gLog.Info("redis connect success")
		gLog.Info("redis connect success")
		gLog.Info("redis connect success")
	}

	return client
}
