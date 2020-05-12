package main

import (
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
	_ "github.com/sagacioushugo/oauth2/generator/generator_default"
	"github.com/sagacioushugo/oauth2/server"
	_ "github.com/sagacioushugo/oauth2/store/store_mem"
	"log"
	"net/http"
)

func main() {
	// modify default config if you need
	config := oauth2.NewDefaultConfig()

	// example: config.ManagerConfig.TokenStoreName = "your token store implemented Mysql, MongoDB, Redis..etc"
	// how to implement your token store refer to store/store_mem/store_mem.go
	config.ManagerConfig.TokenStoreName = "mem"

	// example: config.ManagerConfig.GeneratorName = "your token generator implemented random, jwt, uuid..etc"
	// how to implement your token generator refer to generator/generator_default/generator_default.go
	config.ManagerConfig.GeneratorName = "default"

	oauth2Server := server.NewServer(config)

	// All Oauth2 Pattern must implement
	oauth2Server.SetAuthenticateClientHandler(authenticateClientHandler)

	// PasswordCredentials Pattern must implement
	oauth2Server.SetCheckUserPasswordHandler(checkUserPasswordHandler)

	// Authorization Code Pattern must implement
	oauth2Server.SetCheckUserGrantAccessHandler(checkUserGrantAccessHandler)

	// Other custom api detail refer to server/handler.go
	// CustomizedCheckScopeHandler
	// CustomizedAuthorizeErrHandler
	// CustomizedClientCredentialsUserIdHandler
	// CustomizedTokenExtensionFieldsHandler

	// start token GC
	go oauth2Server.Manager.TokenGC()

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		context := context.NewContext()
		context.Reset(w, r)
		oauth2Server.Authorize(context)
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		context := context.NewContext()
		context.Reset(w, r)
		oauth2Server.Token(context)
	})

	log.Println("Start Listen :8081")

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func checkUserPasswordHandler(username, password string, ctx *context.Context) (userId string, err error) {
	// should check username and password from database
	if username == "admin" && password == "123456" {
		userId = username
		return
	} else {
		return "", oauth2.ErrInvalidUsernameOrPassword
	}
}

func checkUserGrantAccessHandler(req *oauth2.Request, ctx *context.Context) (userId string, err error) {
	// should check user is valid in session
	// if user is valid return userId
	// else return "", oauth2.ErrAccessDenied
	return "admin", nil
}

func authenticateClientHandler(ctx *context.Context, clientIdAndSecret ...string) (redirectUris string, err error) {
	if len(clientIdAndSecret) == 0 || len(clientIdAndSecret) > 2 {
		return "", oauth2.ErrInvalidRequest
	}
	const demoRedirectUris = "http://localhost:8080/test/code_to_token|http://localhost:8080/test/implicit_token"

	if len(clientIdAndSecret) == 1 {
		if clientIdAndSecret[0] == "1" {
			return demoRedirectUris, nil
		} else {
			return "", oauth2.ErrInvalidClient
		}
	}

	if clientIdAndSecret[0] == "1" && clientIdAndSecret[1] == "1he5k5ZUrHFjznxN" {
		return demoRedirectUris, nil
	} else {
		return "", oauth2.ErrInvalidClient
	}

}
