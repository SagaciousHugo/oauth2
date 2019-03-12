package store

import (
	"github.com/astaxie/beego/context"
)

type TokenStore interface {
	Init(tokenConfig string) error
	NewToken(ctx *context.Context) Token
	Create(token Token) error
	GetByAccess(access string) (Token, error)
	GetByRefresh(fresh string) (Token, error)
	GetByCode(code string) (Token, error)
	CreateAndDel(tokenNew Token, tokenDel Token) error
	GC(gcInterval int64)
}
