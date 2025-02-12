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

	var response []model.Activity

	// Get all Comments
	comments, err := c.Service.ItemActivity.GetItemActivityByItemId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, comment := range comments {
		response = append(response, model.Activity{
			Action:  "comment",
			Created: comment.Created,
			Comment: comment,
		})
	}

	// Get legal consultation request
	consultations, err := c.Service.LegalConsultation.GetLegalConsultationByItemId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(consultations) > 0 {
		// get all emails of legal consultants
		var legalConsultants []string
		for _, consultation := range consultations {
			legalConsultants = append(legalConsultants, consultation.Email)
		}

		// join all emails of legal consultants delimited by comma
		legalConsultantsStr := ""
		for i, v := range legalConsultants {
			if i > 0 {
				legalConsultantsStr += ", "
			}
			legalConsultantsStr += v
		}

		response = append(response, model.Activity{
			Action: "consult",
			Actor:  consultations[0].CreatedBy,
			Comment: model.ItemActivity{
				Content: legalConsultantsStr,
			},
			Created: consultations[0].Created,
		})
	}

	// Arrange response by response.Created
	for i := 0; i < len(response); i++ {
		for j := i + 1; j < len(response); j++ {
			if response[i].Created > response[j].Created {
				response[i], response[j] = response[j], response[i]
			}
		}
	}

	// Encode to JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
