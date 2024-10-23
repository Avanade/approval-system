package authenticator

import (
	"main/model"
	"net/http"
)

type AuthenticatorService interface {
	AccessTokenIsValid(r *http.Request) bool
	ClearFromSession(w *http.ResponseWriter, r *http.Request, session string) error
	GetAuthCodeURL(state string) string
	GetAuthenticatedUser(r *http.Request) (*model.AzureUser, error)
	GetLogoutURL() (string, error)
	GetStringValue(r *http.Request, session string, key string) (string, error)
	IsAuthenticated(r *http.Request) (bool, bool, error)
	ProcessToken(code string) (*UserToken, error)
	RefreshToken(w *http.ResponseWriter, r *http.Request) error
	SaveOnSession(w *http.ResponseWriter, r *http.Request, session string, p ...interface{}) error
}
