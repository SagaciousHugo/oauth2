package store_mysql

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/sagacioushugo/oauth2"
	"github.com/sagacioushugo/oauth2/store"
	"time"
)

type TokenStore struct {
}

func (s *TokenStore) Init(tokenConfig string) error {
	return nil
}

func (s *TokenStore) NewToken(ctx *context.Context) store.Token {
	token := Token{
		GrantSessionId: ctx.Input.CruSession.SessionID(),
	}
	return &token
}

func (s *TokenStore) Create(token store.Token) error {
	t, _ := token.(*Token)
	o := orm.NewOrm()
	_, err := o.Insert(t)
	return err
}

func (s *TokenStore) GetByAccess(access string) (store.Token, error) {
	o := orm.NewOrm()
	var token Token

	if err := o.QueryTable("token").Filter("access", access).One(&token); err != nil {
		if err == orm.ErrNoRows {
			return nil, oauth2.ErrInvalidAccessToken
		} else {
			return nil, err
		}
	} else {
		return &token, nil
	}
}
func (s *TokenStore) GetByRefresh(refresh string) (store.Token, error) {
	o := orm.NewOrm()
	var token Token

	if err := o.QueryTable("token").Filter("refresh", refresh).One(&token); err != nil {
		if err == orm.ErrNoRows {
			return nil, oauth2.ErrInvalidRefreshToken
		} else {
			return nil, err
		}
	} else {
		return &token, nil
	}

}
func (s *TokenStore) GetByCode(code string) (store.Token, error) {
	o := orm.NewOrm()
	var token Token

	if err := o.QueryTable("token").Filter("code", code).One(&token); err != nil {
		if err == orm.ErrNoRows {
			return nil, oauth2.ErrInvalidAuthorizeCode
		} else {
			return nil, err
		}
	} else {
		return &token, nil
	}
}
func (s *TokenStore) CreateAndDel(tokenNew store.Token, tokenDel store.Token) error {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return err
	}
	//new, _ := tokenNew.(*Token)
	del, _ := tokenDel.(*Token)

	if _, err := o.QueryTable("token").Filter("id", del.Id).Delete(); err != nil {
		if rerr := o.Rollback(); rerr != nil {
			logs.Error(rerr)
		}
		return err
	} else if _, err := o.Insert(tokenNew); err != nil {
		if rerr := o.Rollback(); rerr != nil {
			logs.Error(rerr)
		}
		return err
	} else {
		if cerr := o.Commit(); cerr != nil {
			logs.Error(cerr)
		}
		return nil
	}
}
func (s *TokenStore) GC(gcInterval int64) {
	timestamp := time.Now().Unix() - gcInterval
	sql := fmt.Sprintf("delete from token where access_create_at is null or unix_timestamp(access_create_at) < %d", timestamp)
	o := orm.NewOrm()

	if _, err := o.Raw(sql).Exec(); err != nil {
		logs.Error(fmt.Errorf("token gc failed: %s", err.Error()))
	}
}
