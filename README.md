# Golang Oauth2 Server
Based on the secondary development `github.com/go-oauth2`, provides TokenStore and TokenGenerator based on interface can be pull out plug type design, can easily customize their need TokenStore and the Generator, and support the TokenGC.


## Features
* clear configuration items with default configuration(config.go)
* 7 custom hooks are supported to extend logic in the authorization process (server/handler.go)
* support response in hook functions (such as output authorization page, error page, etc., subsequent authorization logic is automatically interrupted)
* support for Token GC
* support custom Token Store with any database (you need to implement the store.tokenstore interface yourself)
* support custom Token Generator with any algorithm (you need to implement the Generator.Generator interface yourself)
* based entirely on [RFC 6749](https://tools.ietf.org/html/rfc6749) implementation (support RedirectURI is empty, support multiple RedirectURI configuration)

## Quick Start
### Download and install
```bash
go get github.com/sagacioushugo/oauth2
```

### Build and run example

```bash
go build server.go
./server

go build client.go
./client
```


### Request to client in browser

After receiving the corresponding request, the client will automatically simulate a certain authorization pattern and return the token obtained to the browser

Test the authorization process with [RFC 6749](https://tools.ietf.org/html/rfc6749) and client logs


Authorization Pattern| Client Url | note
---|---|---
Authorization Code|http://localhost:8080/test/authorization_code |
Implicit|http://localhost:8080/test/implicit |
Password Credentials|http://localhost:8080/test/password|
Client Credentials | http://localhost:8080/test/client_credentials|
Refreshing an access token|http://localhost:8080/test/refresh_token?refresh_token=yourtoken| param refresh_token is required


## Documentation
* [English](./README.md)
* [中文文档](./README_CH.md)

## Apache License 2.0

  Copyright (c) 2019 Guoyiming