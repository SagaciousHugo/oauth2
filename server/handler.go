package server

import (
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
	"github.com/sagacioushugo/oauth2/store"
)

type (
	// check user if has or not grant access ( check user is login and user has granted access to the client)
	CheckUserGrantAccessHandler func(req *oauth2.Request, ctx *context.Context) (userId string, err error)

	// check username and password is match
	CheckUserPasswordHandler func(username, password string, ctx *context.Context) (userId string, err error)

	// authenticate client (client is valid or satisfy some condition) or default when authorize request check client existed and token request check client_id match client_secret
	AuthenticateClientHandler func(ctx *context.Context, clientIdAndSecret ...string) (redirectUris string, err error)

	// check scope or default not check scope
	CustomizedCheckScopeHandler func(scope string, grantType *oauth2.GrantType, ctx *context.Context) (allowed bool, err error)

	// check refresh scope or default check newScope equals oldScope
	CustomizedRefreshingScopeHandler func(newScope, oldScope string, ctx *context.Context) (allowed bool, err error)

	CustomizedAuthorizeErrHandler func(err error, ctx *context.Context) error

	CustomizedClientCredentialsUserIdHandler func(clientId string) (userId string)

	CustomizedTokenExtensionFieldsHandler func(token store.Token, req *oauth2.Request, ctx *context.Context) (fieldsValue map[string]interface{})
)
