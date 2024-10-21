package authenticator

import (
	"net/http"

	"golang.org/x/oauth2"
)

type AuthenticatorService interface {
	AccessTokenIsValid(r *http.Request) bool
	GetConfig() *oauth2.Config
}
