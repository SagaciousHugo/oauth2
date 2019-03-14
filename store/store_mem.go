package store

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
	"sync"
	"time"
)

type MemToken struct {
	Id              int64
	ClientId        string
	UserId          string
	Scope           string
	Access          string
	AccessCreateAt  time.Time
	AccessExpireIn  int64
	Refresh         string
	RefreshCreateAt time.Time
	RefreshExpireIn int64
	Code            string
	CodeCreateAt    time.Time
	CodeExpireIn    int64
}

func (token *MemToken) GetClientId() string {
	return token.ClientId
}
func (token *MemToken) SetClientId(clientId string) {
	token.ClientId = clientId
}
func (token *MemToken) GetUserId() string {
	return token.UserId
}
func (token *MemToken) SetUserId(userId string) {
	token.UserId = userId
}
func (token *MemToken) GetScope() string {
	return token.Scope
}
func (token *MemToken) SetScope(scope string) {
	token.Scope = scope
}

//code info
func (token *MemToken) GetCode() string {
	return token.Code
}
func (token *MemToken) SetCode(code string) {
	token.Code = code
}
func (token *MemToken) GetCodeCreateAt() time.Time {
	return token.CodeCreateAt
}
func (token *MemToken) SetCodeCreateAt(codeCreateAt time.Time) {
	token.CodeCreateAt = codeCreateAt
}
func (token *MemToken) GetCodeExpireIn() int64 {
	return token.CodeExpireIn
}
func (token *MemToken) SetCodeExpireIn(codeExpireIn int64) {
	token.CodeExpireIn = codeExpireIn
}

//access info
func (token *MemToken) GetAccess() string {
	return token.Access
}
func (token *MemToken) SetAccess(access string) {
	token.Access = access
}
func (token *MemToken) GetAccessCreateAt() time.Time {
	return token.AccessCreateAt
}
func (token *MemToken) SetAccessCreateAt(accessCreateAt time.Time) {
	token.AccessCreateAt = accessCreateAt
}
func (token *MemToken) GetAccessExpireIn() int64 {
	return token.AccessExpireIn

}
func (token *MemToken) SetAccessExpireIn(accessExpireIn int64) {
	token.AccessExpireIn = accessExpireIn
}

// refresh info
func (token *MemToken) GetRefresh() string {
	return token.Refresh
}
func (token *MemToken) SetRefresh(refresh string) {
	token.Refresh = refresh
}
func (token *MemToken) GetRefreshCreateAt() time.Time {
	return token.RefreshCreateAt
}
func (token *MemToken) SetRefreshCreateAt(refreshCreateAt time.Time) {
	token.RefreshCreateAt = refreshCreateAt
}
func (token *MemToken) GetRefreshExpireIn() int64 {
	return token.RefreshExpireIn
}
func (token *MemToken) SetRefreshExpireIn(refreshExpireIn int64) {
	token.RefreshExpireIn = refreshExpireIn
}

func (token *MemToken) IsCodeExpired() bool {
	if token.CodeExpireIn == 0 {
		return true
	} else {
		ct := time.Now()
		return ct.After(token.CodeCreateAt.Add(time.Second * time.Duration(token.CodeExpireIn)))
	}
}

func (token *MemToken) IsAccessExpired() bool {
	if token.AccessExpireIn == 0 {
		return true
	} else {
		ct := time.Now()
		return ct.After(token.AccessCreateAt.Add(time.Second * time.Duration(token.AccessExpireIn)))
	}
}

func (token *MemToken) IsRefreshExpired() bool {
	if token.RefreshExpireIn == 0 {
		return true
	} else {
		ct := time.Now()
		return ct.After(token.RefreshCreateAt.Add(time.Second * time.Duration(token.RefreshExpireIn)))
	}
}

type MemTokenStore struct {
	autoIncrement int64
	list          []*MemToken
	lock          sync.RWMutex
}

func (s *MemTokenStore) Init(tokenConfig string) error {
	return nil
}

func (s *MemTokenStore) NewToken(ctx *context.Context) Token {
	token := MemToken{}
	return &token
}

func (s *MemTokenStore) Create(token Token) error {
	t := token.(*MemToken)
	s.lock.Lock()
	defer s.lock.Unlock()
	t.Id = s.autoIncrement
	s.autoIncrement++
	s.list = append(s.list, t)
	return nil
}

func (s *MemTokenStore) GetByAccess(access string) (Token, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.list {
		if v.Access == access {
			return v, nil
		}
	}
	return nil, oauth2.ErrInvalidAccessToken

}
func (s *MemTokenStore) GetByRefresh(refresh string) (Token, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.list {
		if v.Refresh == refresh {
			return v, nil
		}
	}
	return nil, oauth2.ErrInvalidRefreshToken

}
func (s *MemTokenStore) GetByCode(code string) (Token, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.list {
		if v.Code == code {
			return v, nil
		}
	}
	return nil, oauth2.ErrInvalidAuthorizeCode
}
func (s *MemTokenStore) CreateAndDel(tokenNew Token, tokenDel Token) error {
	new := tokenNew.(*MemToken)
	del := tokenDel.(*MemToken)
	s.lock.Lock()
	defer s.lock.Unlock()
	var ok bool
	for i, v := range s.list {
		if v.Id == del.Id {
			s.list = append(s.list[:i], s.list[i+1:]...)
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("del token not found")
	}
	new.Id = s.autoIncrement
	s.autoIncrement++

	s.list = append(s.list, new)
	return nil
}
func (s *MemTokenStore) GC(gcInterval int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now()
	newList := make([]*MemToken, len(s.list), cap(s.list))
	for _, v := range s.list {
		if now.Before(v.AccessCreateAt.Add(time.Second * time.Duration(gcInterval))) {
			newList = append(newList, v)
		}
	}
	s.list = newList
}
