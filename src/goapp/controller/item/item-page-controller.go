package item

import (
	"encoding/json"
	"fmt"
	"main/config"
	"main/model"
	"main/pkg/session"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type itemPageController struct {
	*service.Service
	CommunityPortalAppId string
}

func NewItemPageController(s *service.Service, conf config.ConfigManager) ItemPageController {
	return &itemPageController{
		Service:              s,
		CommunityPortalAppId: conf.GetCommunityPortalAppId(),
	}
}

func (c *itemPageController) MyRequests(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}
	user := model.AzureUser{
		Name:  profile["name"].(string),
		Email: profile["preferred_username"].(string),
	}

	application, err := c.Service.Application.GetApplicationById(c.CommunityPortalAppId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(application)
	if err != nil {
		fmt.Println(err)
		return
	}

	t, d := c.Service.Template.UseTemplate("myrequests", r.URL.Path, user, string(b))

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *itemPageController) MyApprovals(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}
	user := model.AzureUser{
		Name:  profile["name"].(string),
		Email: profile["preferred_username"].(string),
	}

	application, err := c.Service.Application.GetApplicationById(c.CommunityPortalAppId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(application)
	if err != nil {
		fmt.Println(err)
		return
	}

	t, d := c.Service.Template.UseTemplate("myapprovals", r.URL.Path, user, string(b))

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *itemPageController) RespondToItem(w http.ResponseWriter, r *http.Request) {
	session, _ := session.Store.Get(r, "auth-session")

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}
	user := model.AzureUser{
		Name:  profile["name"].(string),
		Email: profile["preferred_username"].(string),
	}

	params := mux.Vars(r)

	itemIsAuthorized, err := c.Service.Item.ItemIsAuthorized(
		params["appGuid"],
		params["appModuleGuid"],
		params["itemGuid"],
		user.Email,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !itemIsAuthorized.IsAuthorized {
		t, d := c.Service.Template.UseTemplate("Unauthorized", r.URL.Path, user, nil)
		err = t.Execute(w, d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if itemIsAuthorized.IsApproved != nil {
			var text string
			if itemIsAuthorized.IsApproved.Value {
				text = "approved"
			} else {
				text = "rejected"
			}

			data := RespondePageData{
				Response: text,
			}

			t, d := c.Service.Template.UseTemplate("already-processed", r.URL.Path, user, data)
			err = t.Execute(w, d)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			item, err := c.Service.Item.GetItemById(params["itemGuid"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := RespondePageData{
				ApplicationId:       params["appGuid"],
				ApplicationModuleId: params["appModuleGuid"],
				ItemId:              params["itemGuid"],
				ApproverEmail:       user.Email,
				IsApproved:          params["isApproved"],
				Data:                *item,
				RequireRemarks:      itemIsAuthorized.RequireRemarks,
			}

			t, d := c.Service.Template.UseTemplate("response", r.URL.Path, user, data)
			err = t.Execute(w, d)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
