package user

import (
	"encoding/json"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type userController struct {
	*service.Service
}

func NewUserController(svc *service.Service) UserController {
	return &userController{
		Service: svc,
	}
}

func (c *userController) SearchUserFromActiveDirectory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]

	users, err := c.Service.MsGraph.SearchUsers(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
