package item

import (
	"encoding/json"
	"fmt"
	"main/config"
	"main/model"
	"main/pkg/session"
	"main/service"
	"net/http"
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
