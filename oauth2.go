package oauth2

import (
	"net/http"
)

type Request struct {
	// authorize
	ResponseType ResponseType
	State        string

	// token
	ClientSecret string
	Code         string
	Refresh      string

	//common
	ClientId    string
	ClientInfo  map[string]string
	GrantType   GrantType
	UserId      string
	RedirectUri string
	Scope       string
}

// Response error response
type Response struct {
	Error       error
	ErrorCode   int
	Description string
	Uri         string
	StatusCode  int
	Header      http.Header
}
