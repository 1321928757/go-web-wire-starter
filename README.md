# go-web-wire-starter

## 简介

使用go-wire框架与gin框架搭建的web开发脚手架，有助于web开发者快速开发curd操作，以及学习go-wire框架的工程化实践

## 项目结构

```bash
├─cmd				// 项目根入口
│  └─app			
├─conf				// 配置文件
├─config			// 配置映射模型
├─example			// 示例代码
│  ├─captcha		// 行为验证码vue示例
│  ├─email			// 邮件功能示例
│  └─storage		// 文件存储功能示例
├─internal			// 核心目录
│  ├─compo			// 其他组件
│  ├─dao			// dao数据层
│  ├─domain			// 业务模型
│  ├─handler		// handler层（controller层）
│  ├─mildware		// 中间件
│  ├─model			// 数据模型
│  ├─pkg			// 其他功能类
│  │  ├─error		// 错误码
│  │  ├─request		// 请求参数模型
│  │  └─response	// 请求响应数据
│  └─service		// service层
├─router			// Http路由
├─static			// 静态文件（如web网页）
├─storage			// 项目生成的日志等文件
│  └─logs			// 项目日志
└─util				// 工具
    ├─fileUtil
    ├─hash
    ├─path
    ├─str
    └─validator
```



## 实现的功能

• 全局依赖注入（go-wire）

• 自定义参数校验规则与校验失败响应

• viper配置统一管理

• Jwt令牌生成，校验，续签，黑名单

• 定时任务

• 文件存储（支持本地，七牛云Kodo，阿里云Oss，腾讯云Cos等存储服务，支持扩展）

• Redis分布式锁

• 请求限流器（基于令牌桶算法）

• 邮件服务

• 自定义命令行命令（代码中以数据库迁移migrate命令为示例）



## 后续待添加的功能

• Mq消息队列



## 技术栈

**gin**: https://gin-gonic.com/zh-cn/	"性能极好的HTTP Web 框架"

**go-wire**: https://github.com/google/wire	"谷歌官方的依赖注入代码生成器"

**gorm**: https://gorm.io/	"强大的Go语言ORM库,用于与数据库交互操作"

**cobra**: https://github.com/spf13/cobra	"用于构建强大的命令行接口的库,可以很方便地组织命令行参数和子命令"

**go-playground/validator**: https://github.com/go-playground/validator	"Go 结构和字段验证"

**golang-jwt/jwt/v5**: https://github.com/golang-jwt/jwt	"web网站的身份令牌"

**zap**: https://github.com/uber-go/zap	"高性能日志库,可以用于记录请求日志"

**pflag**: https://github.com/spf13/pflag	"命令行参数和标志解析库,常与cobra一起使用"

**viper**: https://github.com/spf13/viper	"可以从配置文件(json、yaml等)或环境变量中读取配置的库"

**sonyflake**: https://github.com/sony/sonyflake	"分布式ID生成算法实现,用于生成唯一ID"

........



## 运行

- go build

  ```bash
  $ go generate
  $ go build -o ./bin/ ./...
  $ ./bin/app
  ```

  

- go run

  ```bash
  $ go generate
  $ go run cmd/app/main.go cmd/app/wire_gen.go cmd/app/app.go
  ```

  

- make

  ```bash
  $ make generate
  $ make build
  $ ./bin/app
  ```



- 数据库迁移（cobra命令测试）

  ```bash
  $ go generate
  $ go run cmd/app/main.go cmd/app/wire_gen.go cmd/app/app.go
  ```




## 文件存储功能示例

这里我们以项目中的media_service.go的代码为例，需要使用到internal/compo/中的storage.Storage

```go
type MediaService struct {
    conf     *config.Configuration
    log      *zap.Logger
    mediaDao *dao.MediaDao
    storage  *storage.Storage
}

// NewMediaService .
func NewMediaService(conf *config.Configuration, log *zap.Logger, mediaDao *dao.MediaDao,
    s *storage.Storage) *MediaService {
    return &MediaService{conf: conf, log: log, mediaDao: mediaDao, storage: s}
}
```

1.FileDriver()获取对应存储驱动

func (storage *Storage) FileDriver(disk ...string) (StorageDriver, error)

可传入参数disk选择存储驱动，若不填，则使用配置中的默认驱动

```go
// 获取默认驱动（可通过修改Storage.default配置）
disk, err := s.storage.FileDriver()
// 获取腾讯云Cos存储驱动
disk, err := s.storage.FileDriver(storage.Cos)
// 获取七牛云Kodo存储驱动
disk, err := s.storage.FileDriver(storage.KoDo)
// 获取阿里云Oss服务
disk, err := s.storage.FileDriver(storage.Oss)
```



存储驱动方法

```go
	// Put 上传文件
	Put(key string, r io.Reader, dataLength int64) error
	// PutFile 上传本地文件
	PutFile(key string, localFile string) error
	// Get 获取文件
	Get(key string) (io.ReadCloser, error)
	// Rename 重命名文件
	Rename(srcKey string, destKey string) error
	// Copy 复制文件
	Copy(srcKey string, destKey string) error
	// Exists 判断文件是否存在
	Exists(key string) (bool, error)
	// Size 获取文件大小
	Size(key string) (int64, error)
	// Delete 删除文件
	Delete(key string) error
	// Url 获取文件访问URL
	Url(key string) string
```



简单使用示例

```go
	// 读取文件
	file, err := params.Image.Open()
	defer file.Close()
	if err != nil {
		return nil, cErr.BadRequest("上传失败")
	}

	// 获取存储驱动
	disk, err := s.storage.FileDriver(storage.Oss)
	if err != nil {
		return nil, cErr.BadRequest(s.storage.GetDefaultDiskType() + "disk not found")
	}
	// 生成路径和文件名
	key := s.makeFaceDir(params.Business) + "/" + s.HashName(params.Image.Filename)
	// 上传文件到本地（服务器）
	err = disk.Put(key, file, params.Image.Size)
	if err != nil {
		return nil, cErr.BadRequest("上传失败")
	}
```



## 邮件功能示例

1.配置邮箱服务信息：

```yaml
# conf/config.yaml
email: # 邮件配置
  sender_name: kitie # 发件人名称
  sender_email: xxxxx # 发件人邮箱
  sender_password: xxxxxxxx # 发件人密码
  host: smtp.qq.com # smtp服务器地址
  port: 587 # smtp服务器端口
  max_connection: 4 # 最大并发SMTP连接数
  max_timeout: 20  # 最大超时时间（s）
```

2.使用到的对象为email.EmailDriver，我们这里以CaptchaService验证码服务层举例：// 验证码服务层

```go
type CaptchaService struct {
    conf    *config.Configuration
    log     *zap.Logger
    email   *email.EmailDriver //邮件服务驱动
    rdb     *redis.Client
    captcha *compo.CaptchaCompo
}

func NewCaptchaService(conf *config.Configuration, log *zap.Logger,
    email *email.EmailDriver, rdb *redis.Client, captcha *compo.CaptchaCompo) *CaptchaService {
    return &CaptchaService{conf: conf, log: log, email: email, rdb: rdb, captcha: captcha}
}
```

3.发送邮件

```go
// email为目标邮箱地址，title为邮箱标题，content为邮箱内容
err := s.email.SendRegisterMail(email, title, content)
```



## 行为验证码示例

验证码前端界面可参考example/captcha/案例（vue编写）

安装依赖

```
npm install 
```

运行项目

```
npm run serve
```

go-wire依赖注入代码生成（在cmd/app/目录下）

```
wire gen wire.go
```

