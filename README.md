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
│  └─logs
└─util
    ├─fileUtil
    ├─hash
    ├─path
    ├─str
    └─validator
```



## 技术栈

gin：性能极好的HTTP Web 框架

go-wire：goggle官方的依赖注入代码生成工具

gorm：

go-playground/validator

golang-jwt/jwt/v5

zap

pflag

viper

sonyflake

........