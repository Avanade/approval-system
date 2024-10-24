package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionManager struct {
	w           *http.ResponseWriter
	r           *http.Request
	values      map[string]interface{}
	sessionName string
	store       *sessions.FilesystemStore
}
