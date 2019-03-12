package oauth2

/*

  Oauth2 Pattern          responseType(GET /authorize)            grantType(POST /token)
Authorization Code Grant               code                          authorization_code
Implicit Grant					       token                               x
Password Credentials Grant              x                               password
Client Credentials Grant                x                          client_credentials
Refreshing an access token              x                            refresh_token

Refer to rfc6749
https://tools.ietf.org/html/rfc6749#section-4.1
https://tools.ietf.org/html/rfc6749#section-4.2
https://tools.ietf.org/html/rfc6749#section-4.3
https://tools.ietf.org/html/rfc6749#section-4.4
https://tools.ietf.org/html/rfc6749#section-6

*/

type ResponseType string

func (rt ResponseType) IsValid() bool {
	switch rt {
	case Token:
		return true
	case Code:
		return true
	default:
		return false
	}
}

const (
	Token ResponseType = "token"
	Code  ResponseType = "code"
)

type GrantType string

const (
	AuthorizationCode   GrantType = "authorization_code"
	Implicit            GrantType = "__implicit"
	PasswordCredentials GrantType = "password"
	ClientCredentials   GrantType = "client_credentials"
	RefreshToken        GrantType = "refresh_token"
)

func (gt GrantType) IsValid() bool {
	switch gt {
	case AuthorizationCode:
		return true
	case PasswordCredentials:
		return true
	case ClientCredentials:
		return true
	case RefreshToken:
		return true
	default:
		return false
	}
}
