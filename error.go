package oauth2

import (
	"errors"
)

var (
	// oauth2 err https://tools.ietf.org/html/rfc6749#section-5.2
	ErrInvalidRequest          = errors.New("invalid_request")
	ErrUnauthorizedClient      = errors.New("unauthorized_client")
	ErrAccessDenied            = errors.New("access_denied")
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")
	ErrInvalidScope            = errors.New("invalid_scope")
	ErrServerError             = errors.New("server_error")
	ErrTemporarilyUnavailable  = errors.New("temporarily_unavailable")
	ErrInvalidClient           = errors.New("invalid_client")
	ErrInvalidGrant            = errors.New("invalid_grant")
	ErrUnsupportedGrantType    = errors.New("unsupported_grant_type")

	// extra customized errors
	ErrInvalidRedirectURI        = errors.New("invalid_redirect_uri")
	ErrInvalidAuthorizeCode      = errors.New("invalid_authorize_code")
	ErrInvalidAccessToken        = errors.New("invalid_access_token")
	ErrInvalidRefreshToken       = errors.New("invalid_refresh_token")
	ErrExpiredAuthorizeCode      = errors.New("expired_authorize_code")
	ErrExpiredAccessToken        = errors.New("expired_access_token")
	ErrExpiredRefreshToken       = errors.New("expired_refresh_token")
	ErrInvalidUsernameOrPassword = errors.New("invalid_username_or_password")
)

// Descriptions error description
var Descriptions = map[error]string{
	ErrInvalidRequest:          "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
	ErrUnauthorizedClient:      "The client is not authorized to request an authorization code using this method",
	ErrAccessDenied:            "The resource owner or authorization server denied the request",
	ErrUnsupportedResponseType: "The authorization server does not support obtaining an authorization code using this method",
	ErrInvalidScope:            "The requested scope is invalid, unknown, or malformed",
	ErrServerError:             "The authorization server encountered an unexpected condition that prevented it from fulfilling the request",
	ErrTemporarilyUnavailable:  "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server",
	ErrInvalidClient:           "Client authentication failed",
	ErrInvalidGrant:            "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client",
	ErrUnsupportedGrantType:    "The authorization grant type is not supported by the authorization server",

	ErrInvalidRedirectURI:        "The request is missing redirect uri or includes an invalid redirect uri value",
	ErrInvalidAuthorizeCode:      "The request is missing authorize code or includes an invalid authorize code value",
	ErrInvalidAccessToken:        "The request is missing access token or includes an invalid access token value",
	ErrInvalidRefreshToken:       "The request is missing refresh token or includes an invalid refresh token value",
	ErrExpiredAuthorizeCode:      "The request includes an expired authorize code value",
	ErrExpiredAccessToken:        "The request includes an expired access token value",
	ErrExpiredRefreshToken:       "The request includes an expired refresh token value",
	ErrInvalidUsernameOrPassword: "The request includes invalid username or password",
}

// StatusCodes response error HTTP status code
var StatusCodes = map[error]int{
	ErrInvalidRequest:          400,
	ErrUnauthorizedClient:      401,
	ErrAccessDenied:            403,
	ErrUnsupportedResponseType: 401,
	ErrInvalidScope:            400,
	ErrServerError:             500,
	ErrTemporarilyUnavailable:  503,
	ErrInvalidClient:           401,
	ErrInvalidGrant:            401,
	ErrUnsupportedGrantType:    401,

	ErrInvalidRedirectURI:        400,
	ErrInvalidAuthorizeCode:      400,
	ErrInvalidAccessToken:        400,
	ErrInvalidRefreshToken:       400,
	ErrExpiredAuthorizeCode:      400,
	ErrExpiredAccessToken:        400,
	ErrExpiredRefreshToken:       400,
	ErrInvalidUsernameOrPassword: 400,
}
