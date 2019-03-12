package generator

import (
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
)

type Generator interface {
	Code(req *oauth2.Request, ctx *context.Context) (code string, err error)
	Token(req *oauth2.Request, ctx *context.Context, isGenerateRefresh bool) (access, refresh string, err error)
}
