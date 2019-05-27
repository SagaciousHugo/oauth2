package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
	"net/http"
	"net/url"
	"strings"
)

func GetErrorData(err error) (data map[string]interface{}, statusCode int, header http.Header) {
	data = make(map[string]interface{})
	if oe, ok := err.(*oauth2.Oauth2Error); ok {
		data["error"] = oe.Error()
		data["error_description"] = oe.ErrorDescription()
		statusCode = oe.StatusCode()
	} else if sc, ok := err.(oauth2.StatusCoder); ok {
		data["error"] = err.Error()
		statusCode = sc.StatusCode()
	} else {
		data["error"] = err.Error()
		data["error_description"] = oauth2.ErrServerError.ErrorDescription()
		statusCode = oauth2.ErrServerError.StatusCode()
	}
	return
}

func GetRedirectUri(req *oauth2.Request, data map[string]interface{}) (uri string, err error) {
	if req.RedirectUri == "" {
		err = oauth2.ErrInvalidRequest
		return
	}
	u, err := url.Parse(req.RedirectUri)

	if err != nil {
		return
	}
	q := u.Query()
	if req.State != "" {
		q.Set("state", req.State)
	}

	for k, v := range data {
		q.Set(k, fmt.Sprint(v))
	}
	switch req.ResponseType {
	case oauth2.Code:
		u.RawQuery = q.Encode()
	case oauth2.Token:
		u.RawQuery = ""
		u.Fragment, err = url.QueryUnescape(q.Encode())
		if err != nil {
			return
		}
	default:
		err = oauth2.ErrInvalidRequest
		return
	}
	uri = u.String()
	return
}

func Redirect(req *oauth2.Request, data map[string]interface{}, ctx *context.Context) (err error) {
	uri, err := GetRedirectUri(req, data)
	if err != nil {
		return
	}
	ctx.Redirect(http.StatusFound, uri)
	return
}

func RedirectError(req *oauth2.Request, err error, ctx *context.Context) error {
	data, _, _ := GetErrorData(err)
	return Redirect(req, data, ctx)
}

func ResponseErr(err error, ctx *context.Context) error {
	data, statusCode, header := GetErrorData(err)
	return ResponseToken(data, header, ctx, statusCode)
}
func ResponseToken(data map[string]interface{}, header http.Header, ctx *context.Context, statusCode ...int) (err error) {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	ctx.ResponseWriter.Header().Set("Cache-Control", "no-store")
	ctx.ResponseWriter.Header().Set("Pragma", "no-cache")

	for k := range header {
		ctx.ResponseWriter.Header().Set(k, header.Get(k))
	}
	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}
	ctx.ResponseWriter.WriteHeader(status)

	err = json.NewEncoder(ctx.ResponseWriter).Encode(data)
	return
}

func ParseBasicAuth(basicAuth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(basicAuth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(basicAuth[len(prefix):])
	if err != nil {
		return
	}
	s := bytes.IndexByte(c, ':')
	if s < 0 {
		return
	}
	return string(c[:s]), string(c[s+1:]), true
}
