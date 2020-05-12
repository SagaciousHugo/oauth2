# Golang Oauth2 Server
基于`github.com/go-oauth2`的二次开发，提供了TokenStore和TokenGenerator基于接口的可拔插式设计，能够方便的定制自己需要的TokenStore和Generator，并且支持TokenGC功能。


## 特点
* 清晰的配置项，并且预先设置了默认配置（config.go）
* 支持7个自定义钩子，可在授权流程扩展逻辑（server/handler.go）
* 支持在钩子函数中response（例如输出授权页面、错误页面等，后续授权逻辑自动中断）
* 支持Token GC功能
* 支持用任意数据库自定义Token Store（你需要自己实现store.TokenStore接口）
* 支持用任意算法自定义Token Generator（你需要自己实现generator.Generator接口）
* 完全基于[RFC 6749](https://tools.ietf.org/html/rfc6749)实现（支持RedirectURI为空，支持配置多RedirectURI）



## 快速开始
### 下载和安装
```bash
go get github.com/sagacioushugo/oauth2
```

### 编译运行示例代码

```bash
go build server.go
./server

go build client.go
./client
```

### 在浏览器向client请求

client收到相应请求后会自动模拟某种授权方式，并把获取的token返回至浏览器

可参考[RFC 6749](https://tools.ietf.org/html/rfc6749)文档，并结合client日志了解具体的授权流程

授权方式| client url | 备注
---|---|---
Authorization Code|http://localhost:8080/test/authorization_code |
Implicit|http://localhost:8080/test/implicit |
Password Credentials|http://localhost:8080/test/password|
Client Credentials | http://localhost:8080/test/client_credentials|
Refreshing an access token|http://localhost:8080/test/refresh_token?refresh_token=yourtoken| 需要参数refresh_token


## 文档
* [English](./README.md)
* [中文文档](./README_CH.md)

## Apache License 2.0

  Copyright (c) 2019 Guoyiming