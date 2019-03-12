package store_mysql

import "time"

type Token struct {
	Id              int64 `orm:"pk;auto"`
	ClientId        string
	UserId          string
	Scope           string
	Access          string    `orm:"null;index"`
	AccessCreateAt  time.Time `orm:"null"`
	AccessExpireIn  int64     `orm:"default(0)"`
	Refresh         string    `orm:"null;index"`
	RefreshCreateAt time.Time `orm:"null"`
	RefreshExpireIn int64     `orm:"default(0)"`
	Code            string    `orm:"null;index"`
	CodeCreateAt    time.Time `orm:"null"`
	CodeExpireIn    int64     `orm:"default(0)"`
	GrantSessionId  string    `orm:"default();index"`
}

func (token *Token) GetClientId() string {
	return token.ClientId
}
func (token *Token) SetClientId(clientId string) {
	token.ClientId = clientId
}
func (token *Token) GetUserId() string {
	return token.UserId
}
func (token *Token) SetUserId(userId string) {
	token.UserId = userId
}
func (token *Token) GetScope() string {
	return token.Scope
}
func (token *Token) SetScope(scope string) {
	token.Scope = scope
}

//code info
func (token *Token) GetCode() string {
	return token.Code
}
func (token *Token) SetCode(code string) {
	token.Code = code
}
func (token *Token) GetCodeCreateAt() time.Time {
	return token.CodeCreateAt
}
func (token *Token) SetCodeCreateAt(codeCreateAt time.Time) {
	token.CodeCreateAt = codeCreateAt
}
func (token *Token) GetCodeExpireIn() int64 {
	return token.CodeExpireIn
}
func (token *Token) SetCodeExpireIn(codeExpireIn int64) {
	token.CodeExpireIn = codeExpireIn
}

//access info
func (token *Token) GetAccess() string {
	return token.Access
}
func (token *Token) SetAccess(access string) {
	token.Access = access
}
func (token *Token) GetAccessCreateAt() time.Time {
	return token.AccessCreateAt
}
func (token *Token) SetAccessCreateAt(accessCreateAt time.Time) {
	token.AccessCreateAt = accessCreateAt
}
func (token *Token) GetAccessExpireIn() int64 {
	return token.AccessExpireIn

}
func (token *Token) SetAccessExpireIn(accessExpireIn int64) {
	token.AccessExpireIn = accessExpireIn
}

// refresh info
func (token *Token) GetRefresh() string {
	return token.Refresh
}
func (token *Token) SetRefresh(refresh string) {
	token.Refresh = refresh
}
func (token *Token) GetRefreshCreateAt() time.Time {
	return token.RefreshCreateAt
}
func (token *Token) SetRefreshCreateAt(refreshCreateAt time.Time) {
	token.RefreshCreateAt = refreshCreateAt
}
func (token *Token) GetRefreshExpireIn() int64 {
	return token.RefreshExpireIn
}
func (token *Token) SetRefreshExpireIn(refreshExpireIn int64) {
	token.RefreshExpireIn = refreshExpireIn
}

func (token *Token) IsCodeExpired() bool {
	ct := time.Now()
	return token.CodeExpireIn != 0 &&
		ct.After(token.CodeCreateAt.Add(time.Second*time.Duration(token.CodeExpireIn)))
}

func (token *Token) IsAccessExpired() bool {
	ct := time.Now()
	return token.AccessExpireIn != 0 &&
		ct.After(token.AccessCreateAt.Add(time.Second*time.Duration(token.AccessExpireIn)))
}

func (token *Token) IsRefreshExpired() bool {
	ct := time.Now()
	return token.RefreshExpireIn != 0 &&
		ct.After(token.RefreshCreateAt.Add(time.Second*time.Duration(token.RefreshExpireIn)))
}
