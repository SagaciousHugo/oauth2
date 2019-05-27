package oauth2

import (
	"net/http"
)

type StatusCoder interface {
	StatusCode() int
}

type Oauth2Error struct {
	msg         string
	description string
	statusCode  int
}

func (e Oauth2Error) Error() string {
	return e.msg
}

func (e Oauth2Error) StatusCode() int {
	return e.statusCode
}

func (e Oauth2Error) ErrorDescription() string {
	return e.description
}

func NewError(statusCode int, msg string, description string) *Oauth2Error {
	return &Oauth2Error{statusCode: statusCode, msg: msg, description: description}
}

var (
	// oauth2 err https://tools.ietf.org/html/rfc6749#section-5.2
	ErrInvalidRequest          = NewError(http.StatusBadRequest, "invalid_request", "The request is missing a required parameter, includes an unsupported parameter value (other than grant type), repeats a parameter, includes multiple credentials, utilizes more than one mechanism for authenticating the client, or is otherwise malformed.")
	ErrUnauthorizedClient      = NewError(http.StatusUnauthorized, "unauthorized_client", "The authenticated client is not authorized to use this authorization grant type.")
	ErrAccessDenied            = NewError(http.StatusUnauthorized, "access_denied", "The resource owner or authorization server denied the request")
	ErrUnsupportedResponseType = NewError(http.StatusUnauthorized, "unsupported_response_type", "The authorization server does not support obtaining an authorization code or an access token using this method")
	ErrInvalidScope            = NewError(http.StatusBadRequest, "invalid_scope", "The requested scope is invalid, unknown, or malformed")
	ErrServerError             = NewError(http.StatusInternalServerError, "server_error", "The authorization server encountered an unexpected condition that prevented it from fulfilling the request")
	ErrTemporarilyUnavailable  = NewError(http.StatusServiceUnavailable, "temporarily_unavailable", "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server")
	ErrInvalidClient           = NewError(http.StatusBadRequest, "invalid_client", "Client authentication failed")
	ErrInvalidGrant            = NewError(http.StatusUnauthorized, "invalid_grant", "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client")
	ErrUnsupportedGrantType    = NewError(http.StatusUnauthorized, "unsupported_grant_type", "The authorization grant type is not supported by the authorization server")

	// extra customized errors
	ErrInvalidRedirectURI        = NewError(http.StatusBadRequest, "invalid_redirect_uri", "The request is missing redirect uri or includes an invalid redirect uri value")
	ErrInvalidAuthorizeCode      = NewError(http.StatusBadRequest, "invalid_authorize_code", "The request is missing authorize code or includes an invalid authorize code value")
	ErrInvalidAccessToken        = NewError(http.StatusBadRequest, "invalid_access_token", "The request is missing access token or includes an invalid access token value")
	ErrInvalidRefreshToken       = NewError(http.StatusBadRequest, "invalid_refresh_token", "The request is missing refresh token or includes an invalid refresh token value")
	ErrExpiredAuthorizeCode      = NewError(http.StatusBadRequest, "expired_authorize_code", "The request includes an expired authorize code value")
	ErrExpiredAccessToken        = NewError(http.StatusBadRequest, "expired_access_token", "The request includes an expired access token value")
	ErrExpiredRefreshToken       = NewError(http.StatusBadRequest, "expired_refresh_token", "The request includes an expired refresh token value")
	ErrInvalidUsernameOrPassword = NewError(http.StatusBadRequest, "invalid_username_or_password", "The request includes invalid username or password")
)
