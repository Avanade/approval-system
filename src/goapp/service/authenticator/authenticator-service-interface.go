package authenticator

import (
	"net/http"
)

type AuthenticatorService interface {
	AccessTokenIsValid(r *http.Request) bool
}
