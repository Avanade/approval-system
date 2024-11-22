package item

import (
	"encoding/json"
	"fmt"
	"main/config"
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
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	t, d := c.Service.Template.UseTemplate("myrequests", r.URL.Path, *user, string(b))

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *itemPageController) MyApprovals(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	t, d := c.Service.Template.UseTemplate("myapprovals", r.URL.Path, *user, string(b))

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *itemPageController) RespondToItem(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		t, d := c.Service.Template.UseTemplate("Unauthorized", r.URL.Path, *user, nil)
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

			t, d := c.Service.Template.UseTemplate("already-processed", r.URL.Path, *user, data)
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

			t, d := c.Service.Template.UseTemplate("response", r.URL.Path, *user, data)
			err = t.Execute(w, d)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func (c *itemPageController) ReassignApproval(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		t, d := c.Service.Template.UseTemplate("Unauthorized", r.URL.Path, *user, nil)
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

			t, d := c.Service.Template.UseTemplate("already-processed", r.URL.Path, *user, data)
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
				Data:                *item,
				ApproveText:         params["ApproveText"],
				RejectText:          params["RejectText"],
			}

			t, d := c.Service.Template.UseTemplate("reassign", r.URL.Path, *user, data)
			err = t.Execute(w, d)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func (c *itemPageController) MultipleApprovals(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, d := c.Service.Template.UseTemplate("multiple-approvals", r.URL.Path, *user, nil)

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
