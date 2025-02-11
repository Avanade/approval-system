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

	// Remove user from involvedUserr.Approvers if user is an approver
	for i, v := range involvedUsers.Approvers {
		if v == user.Email {
			involvedUsers.Approvers = append(involvedUsers.Approvers[:i], involvedUsers.Approvers[i+1:]...)
			break
		}
	}

	// Send email notification to approvers
	if len(involvedUsers.Approvers) > 0 {
		err = c.Service.Email.SendActivityEmail(&req, involvedUsers.Approvers, c.Config.GetHomeURL(), "response")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Send email notification to requestor if user is not the requestor
	if involvedUsers.Requestor != user.Email {
		err = c.Service.Email.SendActivityEmail(&req, []string{involvedUsers.Requestor}, c.Config.GetHomeURL(), "view")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
