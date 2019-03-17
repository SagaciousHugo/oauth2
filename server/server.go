package server

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/sagacioushugo/oauth2"
	"github.com/sagacioushugo/oauth2/manager"
	"strings"
	"time"
)

func NewDefaultServer() *Server {
	config := oauth2.NewDefaultConfig()
	m := manager.NewManager(&config.ManagerConfig)
	go m.TokenGC()
	return &Server{
		config:  config,
		Manager: m,
	}
}

func NewServer(config *oauth2.Config) *Server {
	if config == nil {
		panic("oauth2 config must not be nil")
	}
	m := manager.NewManager(&config.ManagerConfig)

	return &Server{
		config:  config,
		Manager: m,
	}

}

type Server struct {
	config  *oauth2.Config
	Manager *manager.Manager

	checkUserGrantAccessHandler CheckUserGrantAccessHandler

	checkUserPasswordHandler CheckUserPasswordHandler

	authenticateClientHandler AuthenticateClientHandler

	customizedCheckScopeHandler CustomizedCheckScopeHandler

	customizedRefreshingScopeHandler CustomizedRefreshingScopeHandler

	customizedAuthorizeErrHandler CustomizedAuthorizeErrHandler

	customizedClientCredentialsUserIdHandler CustomizedClientCredentialsUserIdHandler

	customizedTokenExtensionFieldsHandler CustomizedTokenExtensionFieldsHandler
}

/**
handle authorize
*/
func (s *Server) Authorize(ctx *context.Context) {
	var resErr error
	var tokenData map[string]interface{}
	var req *oauth2.Request
	var userId string
	if req, resErr = s.validationAuthorizeRequest(ctx); resErr != nil {
		goto Resp
	}

	if h := s.checkUserGrantAccessHandler; h != nil {
		if userId, resErr = h(req, ctx); resErr != nil {
			goto Resp
		} else {
			req.UserId = userId
		}
	} else {
		resErr = fmt.Errorf("checkUserGrantAccessHandler is nil")
		goto Resp
	}
	tokenData, resErr = s.generateToken(req, ctx)

Resp:
	if ctx.ResponseWriter.Started {
		return
	}
	if resErr != nil {
		if h := s.customizedAuthorizeErrHandler; h != nil {
			if herr := h(resErr, ctx); herr != nil {
				resErr = herr
			}
		}
		if ctx.ResponseWriter.Started {
			return
		}
		if err := RedirectError(req, resErr, ctx); err != nil {
			ResponseErr(err, ctx)
		}

	} else {
		if req.State != "" {
			tokenData["state"] = req.State
		}
		if err := Redirect(req, tokenData, ctx); err != nil {
			ResponseErr(err, ctx)
		}
	}
}

func (s *Server) validationAuthorizeRequest(ctx *context.Context) (req *oauth2.Request, err error) {
	var redirectUri = ctx.Input.Query("redirect_uri")
	var clientId = ctx.Input.Query("client_id")
	var state = ctx.Input.Query("state")
	var scope = ctx.Input.Query("scope")
	var responseType = oauth2.ResponseType(ctx.Request.FormValue("response_type"))
	var grantType oauth2.GrantType = ""
	var clientRedirectUris string

	req = &oauth2.Request{
		RedirectUri:  redirectUri,
		ResponseType: responseType,
		ClientId:     clientId,
		State:        state,
		Scope:        scope,
	}
	// check invalid request
	if !(ctx.Request.Method == "GET" || ctx.Request.Method == "POST") || clientId == "" {
		err = oauth2.ErrInvalidRequest
		return
	}

	// check unsupported response type
	if !responseType.IsValid() {
		err = oauth2.ErrUnsupportedResponseType
		return
	} else if responseType == oauth2.Code {
		grantType = oauth2.AuthorizationCode

	} else if responseType == oauth2.Token {
		grantType = oauth2.Implicit
	}
	if _, ok := s.config.AllowGrantType[grantType]; !ok {
		err = oauth2.ErrUnsupportedResponseType
		return
	}
	req.GrantType = grantType

	// check invalid client
	if h := s.authenticateClientHandler; h != nil {
		var herr error
		if clientRedirectUris, herr = h(ctx, clientId); herr != nil {
			//err type should be oauth2.ErrInvalidClient or oauth2.ErrUnauthorizedClient
			err = herr
			return
		}
	} else {
		err = fmt.Errorf("authenticateClientHandler is nil")
		return
	}

	// check redirect uri
	if redirectUri == "" && !s.config.RedirectAllowEmpty {
		err = oauth2.ErrInvalidRedirectURI
		return
	} else {
		uris := strings.Split(clientRedirectUris, s.config.RedirectUriSep)
		if redirectUri == "" && len(uris) > 0 {
			redirectUri = uris[0]
		} else if redirectUri != "" {
			match := false
			for _, v := range uris {
				if v == redirectUri {
					match = true
					break
				}
			}
			if !match {
				err = oauth2.ErrInvalidRedirectURI
				return
			}
		} else {
			err = oauth2.ErrInvalidRedirectURI
			return
		}
	}

	// check invalid scope
	if h := s.customizedCheckScopeHandler; h != nil {
		if allowed, herr := h(scope, &grantType, ctx); herr != nil {
			err = herr
			return
		} else if !allowed {
			err = oauth2.ErrInvalidScope
			return
		}
	}

	return
}

/**
handle token
*/
func (s *Server) Token(ctx *context.Context) {
	var req *oauth2.Request
	var resErr error
	var tokenData map[string]interface{}
	if req, resErr = s.validationTokenRequest(ctx); resErr != nil {
		goto Resp
	}
	tokenData, resErr = s.generateToken(req, ctx)

Resp:
	if ctx.ResponseWriter.Started {
		return
	}
	if resErr != nil {
		if err := ResponseErr(resErr, ctx); err != nil {
			logs.Error(err)
		}
	} else {

		if err := ResponseToken(tokenData, nil, ctx); err != nil {
			logs.Error(err)
		}
	}

}

func (s *Server) validationTokenRequest(ctx *context.Context) (req *oauth2.Request, err error) {

	clientId, clientSecret, err := GetClientInfoFromBasicAuth(ctx)
	if err != nil {
		return nil, err
	}
	var code = ctx.Input.Query("code")
	var scope = ctx.Input.Query("scope")
	var refresh = ctx.Input.Query("refresh_token")
	var username = ctx.Input.Query("username")
	var password = ctx.Input.Query("password")
	var userId string
	var grantType = oauth2.GrantType(ctx.Input.Query("grant_type"))
	// check invalid request
	if !(ctx.Request.Method == "GET" || ctx.Request.Method == "POST") || clientId == "" {
		return nil, oauth2.ErrInvalidRequest
	}

	// check unsupported grant type
	if !grantType.IsValid() {
		return nil, oauth2.ErrUnsupportedGrantType
	} else if _, ok := s.config.AllowGrantType[grantType]; !ok {
		return nil, oauth2.ErrUnsupportedGrantType
	}

	// check invalid client
	if h := s.authenticateClientHandler; h != nil {
		if _, herr := h(ctx, clientId, clientSecret); herr != nil {
			return nil, herr
		}
	} else {
		return nil, fmt.Errorf("authenticateClientHandler is nil")
	}

	// check invalid scope
	if grantType != oauth2.AuthorizationCode && grantType != oauth2.RefreshToken {
		if h := s.customizedCheckScopeHandler; h != nil {
			if allowed, herr := h(scope, &grantType, ctx); herr != nil {
				return nil, herr
			} else if !allowed {
				return nil, oauth2.ErrInvalidScope
			}
		}
	}
	// AuthorizationCode and RefreshToken userId consists in their code or token
	switch grantType {
	case oauth2.AuthorizationCode:
		if code == "" {
			return nil, oauth2.ErrInvalidAuthorizeCode
		}
	case oauth2.PasswordCredentials:
		if h := s.checkUserPasswordHandler; h != nil {
			var herr error
			if userId, herr = h(username, password, ctx); herr != nil {
				return nil, herr
			}
		} else {
			return nil, fmt.Errorf("checkUserPasswordHandler is nil")
		}
	case oauth2.ClientCredentials:
		if h := s.customizedClientCredentialsUserIdHandler; h != nil {
			userId = h(clientId)
		} else {
			userId = fmt.Sprintf("client-%s", clientId)
		}
	case oauth2.RefreshToken:
		if refresh == "" {
			return nil, oauth2.ErrInvalidRefreshToken
		}
	}

	req = &oauth2.Request{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    grantType,
		Code:         code,
		UserId:       userId,
		Scope:        scope,
		Refresh:      refresh,
	}
	return req, nil

}

func (s *Server) generateToken(req *oauth2.Request, ctx *context.Context) (map[string]interface{}, error) {

	config := s.config.AllowGrantType[req.GrantType]

	token := s.Manager.TokenNew(ctx)

	now := time.Now()

	switch req.GrantType {
	case oauth2.AuthorizationCode:
		if req.ResponseType == oauth2.Code {
			token.SetUserId(req.UserId)
			token.SetClientId(req.ClientId)
			token.SetScope(req.Scope)
			token.SetCodeExpireIn(config.CodeExpire)
			token.SetCodeCreateAt(now)
			if err := s.Manager.GenerateCode(token, req, ctx); err != nil {
				return nil, err
			}
		} else {
			if code, err := s.Manager.TokenGetByCode(req.Code); err != nil {
				return nil, err
			} else if code == nil || code.GetClientId() != req.ClientId {
				return nil, oauth2.ErrInvalidAuthorizeCode
			} else if code.IsCodeExpired() {
				return nil, oauth2.ErrExpiredAuthorizeCode
			} else {
				token.SetUserId(code.GetUserId())
				token.SetClientId(code.GetClientId())
				token.SetScope(code.GetScope())
				token.SetAccessExpireIn(config.AccessTokenExpire)
				token.SetAccessCreateAt(now)
				if config.IsGenerateRefresh {
					token.SetRefreshExpireIn(config.RefreshTokenExpire)
					token.SetRefreshCreateAt(now)
				}
				if err := s.Manager.GenerateTokenAndDelToken(token, code, req, ctx, config.IsGenerateRefresh); err != nil {
					return nil, err
				}
			}

		}
	case oauth2.RefreshToken:
		if refresh, err := s.Manager.TokenGetByRefresh(req.Refresh); err != nil {
			return nil, err
		} else if refresh == nil || refresh.GetClientId() != req.ClientId {
			return nil, oauth2.ErrInvalidRefreshToken
		} else if refresh.IsRefreshExpired() {
			return nil, oauth2.ErrExpiredRefreshToken
		} else {
			if h := s.customizedRefreshingScopeHandler; h != nil {
				if allowed, herr := h(req.Scope, refresh.GetScope(), ctx); herr != nil {
					return nil, herr
				} else if !allowed {
					return nil, oauth2.ErrInvalidScope
				}
				token.SetScope(req.Scope)
			} else {
				token.SetScope(refresh.GetScope())
			}
			token.SetUserId(refresh.GetUserId())
			token.SetClientId(refresh.GetClientId())
			token.SetAccessExpireIn(config.AccessTokenExpire)
			token.SetAccessCreateAt(now)
			if config.IsGenerateRefresh {
				token.SetRefreshExpireIn(config.RefreshTokenExpire)
				if config.IsResetRefreshTime {
					token.SetRefreshCreateAt(now)
				} else {
					token.SetRefreshCreateAt(refresh.GetRefreshCreateAt())

				}
			}
			if err := s.Manager.GenerateTokenAndDelToken(token, refresh, req, ctx, config.IsGenerateRefresh); err != nil {
				return nil, err
			}
		}
	default:
		// Implicit PasswordCredentials ClientCredentials
		token.SetScope(req.Scope)
		token.SetUserId(req.UserId)
		token.SetClientId(req.ClientId)
		token.SetAccessExpireIn(config.AccessTokenExpire)
		token.SetAccessCreateAt(now)
		if config.IsGenerateRefresh {
			token.SetRefreshExpireIn(config.RefreshTokenExpire)
			token.SetRefreshCreateAt(now)
		}
		if err := s.Manager.GenerateToken(token, req, ctx, config.IsGenerateRefresh); err != nil {
			return nil, err
		}
	}
	data := make(map[string]interface{})

	if req.GrantType == oauth2.AuthorizationCode && req.ResponseType == oauth2.Code {
		data["code"] = token.GetCode()
	} else {
		data["access_token"] = token.GetAccess()
		data["token_type"] = s.config.TokenType
		data["expires_in"] = token.GetAccessExpireIn()

		if config.IsGenerateRefresh {
			data["refresh_token"] = token.GetRefresh()
		}

		if h := s.customizedTokenExtensionFieldsHandler; h != nil {
			ext := h(token, req, ctx)
			for k, v := range ext {
				if _, ok := data[k]; ok {
					continue
				}
				data[k] = v
			}
		}
	}

	return data, nil
}

func (s *Server) SetCheckUserGrantAccessHandler(h CheckUserGrantAccessHandler) {
	s.checkUserGrantAccessHandler = h
}

func (s *Server) SetCheckUserPasswordHandler(h CheckUserPasswordHandler) {
	s.checkUserPasswordHandler = h
}

func (s *Server) SetAuthenticateClientHandler(h AuthenticateClientHandler) {
	s.authenticateClientHandler = h
}

func (s *Server) SetCustomizedCheckScopeHandler(h CustomizedCheckScopeHandler) {
	s.customizedCheckScopeHandler = h
}

func (s *Server) SetCustomizedRefreshingScopeHandler(h CustomizedRefreshingScopeHandler) {
	s.customizedRefreshingScopeHandler = h
}

func (s *Server) SetCustomizedAuthorizeErrHandler(h CustomizedAuthorizeErrHandler) {
	s.customizedAuthorizeErrHandler = h
}

func (s *Server) SetCustomizedClientCredentialsUserIdHandler(h CustomizedClientCredentialsUserIdHandler) {
	s.customizedClientCredentialsUserIdHandler = h
}

func (s *Server) SetCustomizedTokenExtensionFieldsHandler(h CustomizedTokenExtensionFieldsHandler) {
	s.customizedTokenExtensionFieldsHandler = h
}
