package store

import (
	"time"
)

type Token interface {

	//base info
	GetClientId() string
	SetClientId(clientId string)
	GetUserId() string
	SetUserId(userId string)
	GetScope() string
	SetScope(scope string)

	//code info
	GetCode() string
	SetCode(code string)
	GetCodeCreateAt() time.Time
	SetCodeCreateAt(codeCreateAt time.Time)
	GetCodeExpireIn() int64
	SetCodeExpireIn(codeExpireIn int64)

	//access info
	GetAccess() string
	SetAccess(access string)
	GetAccessCreateAt() time.Time
	SetAccessCreateAt(accessCreateAt time.Time)
	GetAccessExpireIn() int64
	SetAccessExpireIn(accessExpireIn int64)

	// refresh info
	GetRefresh() string
	SetRefresh(refresh string)
	GetRefreshCreateAt() time.Time
	SetRefreshCreateAt(refreshCreateAt time.Time)
	GetRefreshExpireIn() int64
	SetRefreshExpireIn(refreshExpireIn int64)

	IsAccessExpired() bool
	IsCodeExpired() bool
	IsRefreshExpired() bool
}
