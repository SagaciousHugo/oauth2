package server

import (
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
	resp := new(oauth2.Response)
	if v, ok := oauth2.Descriptions[err]; ok {
		resp.Error = err
		resp.Description = v
		resp.StatusCode = oauth2.StatusCodes[err]
	} else {
		resp.Error = oauth2.ErrServerError
		resp.Description = err.Error()
		resp.StatusCode = oauth2.StatusCodes[oauth2.ErrServerError]

	}

	data = make(map[string]interface{})

	if err := resp.Error; err != nil {
		data["error"] = err.Error()
	}

	if v := resp.ErrorCode; v != 0 {
		data["error_code"] = v
	}

	if v := resp.Description; v != "" {
		data["error_description"] = v
	}

	if v := resp.Uri; v != "" {
		data["error_uri"] = v
	}

	header = resp.Header

	statusCode = http.StatusInternalServerError
	if v := resp.StatusCode; v > 0 {
		statusCode = v
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

func GetClientInfoFromBasicAuth(ctx *context.Context) (clientId, clientSecret string, err error) {
	header := ctx.Input.Header("Authorization")
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", oauth2.ErrInvalidClient
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", err
	}
	authStr := strings.Split(string(decodeBytes), ":")
	if len(authStr) != 2 {
		return "", "", oauth2.ErrInvalidClient
	} else {
		if clientId, err = url.QueryUnescape(authStr[0]); err != nil {
			return "", "", err
		} else if clientSecret, err = url.QueryUnescape(authStr[1]); err != nil {
			return "", "", err
		} else {
			return clientId, clientSecret, nil
		}
	}

}
