package session

import (
	"net/http"
)

type ConnectSession interface {
	ConnectSession(w *http.ResponseWriter, r *http.Request, session string) SessionManager
}

type SessionManager interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{})
	Save(maxAge int) error
	Destroy() error
}
