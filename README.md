# go-web-wire-starter

## 简介

使用go-wire框架与gin框架搭建的web开发脚手架，有助于web开发者快速开发curd操作，以及学习go-wire框架的工程化实践

## 项目结构

```bash
├─cmd				// 项目根入口
│  └─app			
├─conf				// 配置文件
├─config			// 配置映射模型
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



## 技术栈

[gin](https://gin-gonic.com/zh-cn/ "性能极好的HTTP Web 框架")
[go-wire](https://github.com/google/wire "谷歌官方的依赖注入代码生成器")
[gorm]: https://gorm.io/	"强大的Go语言ORM库,用于与数据库交互操作"
[cobra]: https://github.com/spf13/cobra	"用于构建强大的命令行接口的库,可以很方便地组织命令行参数和子命令"
[go-playground/validator]: https://github.com/go-playground/validator	"Go 结构和字段验证"
[golang-jwt/jwt/v5]: https://github.com/golang-jwt/jwt	"web网站的身份令牌"
[zap]: https://github.com/uber-go/zap	"高性能日志库,可以用于记录请求日志"
[pflag]: https://github.com/spf13/pflag	"命令行参数和标志解析库,常与cobra一起使用"
[viper]: https://github.com/spf13/viper	"可以从配置文件(json、yaml等)或环境变量中读取配置的库"
[sonyflake]: https://github.com/sony/sonyflake	"分布式ID生成算法实现,用于生成唯一ID"

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

  