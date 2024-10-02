package user

import "net/http"

type UserController interface {
	SearchUserFromActiveDirectory(w http.ResponseWriter, r *http.Request)
}
