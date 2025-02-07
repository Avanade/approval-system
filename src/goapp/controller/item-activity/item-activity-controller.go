package itemActivity

import (
	"encoding/json"
	"main/config"
	"main/model"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type itemActivityController struct {
	Service *service.Service
	Config  config.ConfigManager
}

func NewItemActivityController(s *service.Service, c config.ConfigManager) ItemActivityController {
	return &itemActivityController{
		Service: s,
		Config:  c}
}

func (c *itemActivityController) InsertItemActivity(w http.ResponseWriter, r *http.Request) {
	// Decode payload
	var req model.ItemActivity
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user info
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.CreatedBy = user.Email

	// Insert Item Activity
	_, err = c.Service.ItemActivity.AddItemActivity(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get Involved Users
	involvedUsers, err := c.Service.Item.GetInvolvedUsers(req.ItemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove user from involvedUsers
	for i, v := range involvedUsers {
		if v == user.Email {
			involvedUsers = append(involvedUsers[:i], involvedUsers[i+1:]...)
			break
		}
	}

	// Send email notification
	err = c.Service.Email.SendActivityEmail(&req, involvedUsers, c.Config.GetHomeURL())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *itemActivityController) GetItemActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get all Item Activity
	itemActivity, err := c.Service.ItemActivity.GetItemActivityByItemId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode to JSON
	err = json.NewEncoder(w).Encode(itemActivity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
