package item

import (
	"encoding/json"
	"fmt"
	"main/model"
	"main/pkg/session"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type itemController struct {
	*service.Service
}

func NewItemController(svc *service.Service) ItemController {
	return &itemController{
		Service: svc,
	}
}

// GetItems is a function to get all items
func (c *itemController) GetItems(w http.ResponseWriter, r *http.Request) {
	// Get all items
	var itemOptions model.ItemOptions

	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	user := fmt.Sprintf("%s", profile["preferred_username"])
	itemOptions.User = user

	vars := mux.Vars(r)

	itemOptions.ItemType, err = strconv.ParseInt(vars["type"], 10, 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	itemOptions.ItemStatus, err = strconv.ParseInt(vars["status"], 10, 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()

	if params.Has("offset") {
		itemOptions.Offset, err = strconv.Atoi(params["offset"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		itemOptions.Offset = 0
	}

	if params.Has("filter") {
		itemOptions.Filter, err = strconv.Atoi(params["filter"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		itemOptions.Filter = 10
	}

	if params.Has("search") {
		itemOptions.Search = params["search"][0]
	} else {
		itemOptions.Search = ""
	}

	if params.Has("requestType") {
		itemOptions.RequestType = params["requestType"][0]
	} else {
		itemOptions.RequestType = ""
	}

	if params.Has("organization") {
		itemOptions.Organization = params["organization"][0]
	} else {
		itemOptions.Organization = ""
	}

	result, err := c.Service.Item.GetAll(itemOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (c *itemController) CreateItem(w http.ResponseWriter, r *http.Request) {
	// Decode payload
	var req model.TypItemInsert
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get application module
	appModule, err := c.Service.ApplicationModule.GetApplicationModuleByIdAndApplicationId(req.ApplicationId, req.ApplicationModuleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if application module is valid
	if appModule == nil {
		http.Error(w, "Application module not found", http.StatusBadRequest)
		return
	}

	// Add item to database
	id, err := c.Service.Item.InsertItem(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert approvers to database
	for _, to := range req.Emails {
		err := c.Service.Item.InsertApprovalRequestApprover(model.ApprovalRequestApprover{
			ItemId:        id,
			ApproverEmail: to,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Send email
	err = c.Service.Email.SendApprovalRequestEmail(&req, appModule, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.Service.Item.UpdateItemDateSent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["itemId"] = fmt.Sprintf("%v", id)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)

}
