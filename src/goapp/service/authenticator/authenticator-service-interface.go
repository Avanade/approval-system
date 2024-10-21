package authenticator

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

type AuthenticatorService interface {
	VerifyAccessToken(r *http.Request) (*jwt.Token, error)
}
