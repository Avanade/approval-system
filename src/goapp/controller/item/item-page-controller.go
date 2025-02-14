package item

import (
	"encoding/json"
	"fmt"
	"main/config"
	"main/model"
	"main/service"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type itemPageController struct {
	*service.Service
	CommunityPortalAppId string
	IPDRModuleId         string
	CTO                  string
}

func NewItemPageController(s *service.Service, conf config.ConfigManager) ItemPageController {
	return &itemPageController{
		Service:              s,
		CommunityPortalAppId: conf.GetCommunityPortalAppId(),
		IPDRModuleId:         conf.GetIPDRModuleId(),
		CTO:                  conf.GetCTO(),
	}
}

func (c *itemPageController) ForReview(w http.ResponseWriter, r *http.Request) {
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

	t, d := c.Service.Template.UseTemplate("forreview", r.URL.Path, *user, string(b))

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	// Get the authenticated user
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the parameters from the URL
	params := mux.Vars(r)

	action := params["action"]
	isApprover := action == "response"

	isApproved := params["isApproved"]
	if isApproved == "" {
		isApproved = "true"
	}

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

	var approverResponse string
	if itemIsAuthorized.IsApproved == nil {
		itemIsAuthorized.IsApproved = &model.NullBool{Value: false}
	} else {
		if itemIsAuthorized.IsApproved.Value {
			approverResponse = "approved"
		} else {
			approverResponse = "rejected"
		}
	}
	item, err := c.Service.Item.GetItemById(params["itemGuid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get involved users
	involvedUsers, err := c.Service.Item.GetInvolvedUsers(params["itemGuid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is one among the consultants
	isConsultant := false
	for _, consultant := range involvedUsers.Consultants {
		if user.Email == consultant {
			isConsultant = true
			break
		}
	}

	// If the user is not the approver nor the requestor, show unauthorized page
	if (!itemIsAuthorized.IsAuthorized && isApprover) || (item.RequestedBy != user.Email && action == "view") || (!isConsultant && action == "review") {
		t, d := c.Service.Template.UseTemplate("Unauthorized", r.URL.Path, *user, nil)
		err = t.Execute(w, d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if the item is an IP Disclosure Request
	consultLegalButton := strings.EqualFold(item.ModuleId, c.IPDRModuleId)

	// Check if the user is the CTO
	if consultLegalButton && user.Email != c.CTO {
		consultLegalButton = false
	}

	// Check if legal has been involved
	if consultLegalButton {
		consultations, err := c.Service.LegalConsultation.GetLegalConsultationByItemId(params["itemGuid"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(consultations) > 0 {
			consultLegalButton = false
		}
	}

	// If the user is the approver or the requestor, show the response page
	data := RespondePageData{
		ApplicationId:       params["appGuid"],
		ApplicationModuleId: params["appModuleGuid"],
		ItemId:              params["itemGuid"],
		ApproverEmail:       user.Email,
		IsApproved:          isApproved,
		Data:                *item,
		RequireRemarks:      itemIsAuthorized.RequireRemarks,
		IsApprover:          isApprover, // show remarks and approve/reject buttons
		AlreadyProcessed:    itemIsAuthorized.IsApproved.Value,
		ApproverResponse:    approverResponse,
		ConsultLegalButton:  consultLegalButton,
		Action:              action,
	}

	t, d := c.Service.Template.UseTemplate("response", r.URL.Path, *user, data)
	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	t, d := c.Service.Template.UseTemplate("multiple-approvals", r.URL.Path, *user, string(b))

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
