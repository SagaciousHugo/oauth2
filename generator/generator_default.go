package generator

import (
	"bytes"
	"encoding/base64"
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
)

type Default struct {
}

func (g *Default) Code(req *oauth2.Request, ctx *context.Context) (code string, err error) {
	buf := bytes.NewBufferString(req.ClientId)
	buf.WriteString(req.UserId)
	buf.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	token := uuid.NewV3(uuid.Must(uuid.NewV1(), nil), buf.String())
	code = base64.URLEncoding.EncodeToString(token.Bytes())
	code = strings.TrimRight(code, "=")

	return
}

func (g *Default) Token(req *oauth2.Request, ctx *context.Context, isGenRefresh bool) (access, refresh string, err error) {
	buf := bytes.NewBufferString(req.ClientId)
	buf.WriteString(req.UserId)
	buf.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))

	access = base64.URLEncoding.EncodeToString(uuid.NewV3(uuid.Must(uuid.NewV4(), nil), buf.String()).Bytes())
	access = strings.TrimRight(access, "=")
	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString(uuid.NewV5(uuid.Must(uuid.NewV4(), nil), buf.String()).Bytes())
		refresh = strings.TrimRight(refresh, "=")
	}

	return
}
