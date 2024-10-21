package authenticator

import (
	"net/http"
)

type AuthenticatorService interface {
	AccessTokenIsValid(r *http.Request) bool
	GetAuthCodeURL(state string) string
	GetLogoutURL() (string, error)
	ProcessToken(code string) (*UserToken, error)
}
