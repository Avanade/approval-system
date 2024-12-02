package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
	"sync"
	"time"

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

	user, err := c.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	itemOptions.User = user.Email

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

func (c *itemController) GetItemsByApprover(w http.ResponseWriter, r *http.Request) {
	// Get all items
	var filterOptions model.FilterOptions

	user, err := c.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()

	var requestType string
	if params.Has("requestType") {
		requestType = params["requestType"][0]
	}

	var organization string
	if params.Has("organization") {
		organization = params["organization"][0]
	}

	filterOptions.Page = 0 // Default page is 1 which is 0 in the database
	if params.Has("page") {
		filterOptions.Page, err = strconv.Atoi(params["page"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filterOptions.Page = filterOptions.Page - 1
	}

	if params.Has("filter") {
		filterOptions.Filter, err = strconv.Atoi(params["filter"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		filterOptions.Filter = 50
	}

	result, total, err := c.Service.Item.GetByApprover(user.Email, requestType, organization, filterOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response GetItemsByApproverResponse
	for _, item := range result {
		itemResponse := Item{
			Id:            item.Id,
			Subject:       item.Subject,
			Application:   item.Application,
			ApplicationId: item.ApplicationId,
			Module:        item.Module,
			ModuleId:      item.ModuleId,
			RequestedBy:   item.RequestedBy,
			RequestedOn:   item.Created,
			Approvers:     item.Approvers,
			Body:          item.Body,
		}

		response.Data = append(response.Data, itemResponse)
	}
	if len(response.Data) == 0 {
		response.Data = []Item{}
	}
	response.Page = filterOptions.Page + 1
	response.Filter = filterOptions.Filter
	response.Total = total

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *itemController) CreateItem(w http.ResponseWriter, r *http.Request) {
	// Decode payload
	var req model.ItemInsertRequest
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
		err := c.Service.ApprovalRequestApprover.InsertApprovalRequestApprover(model.ApprovalRequestApprover{
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
		fmt.Println("Error sending email: ", err)
		err = nil
	} else {
		err = c.Service.Item.UpdateItemDateSent(id)
		if err != nil {
			fmt.Println("Error updating DateSent column ", err)
			err = nil
		}
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

func (c *itemController) ProcessResponse(w http.ResponseWriter, r *http.Request) {
	// Decode payload
	var req model.ProcessResponseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate payload
	valid, err := c.Service.Item.ValidateItem(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update item response
	err = c.Service.Item.UpdateItemResponse(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Post callback
	c.postCallback(req.ItemId)

	// Prepare response
	w.WriteHeader(http.StatusOK)
}

func (c *itemController) ProcessMultipleResponse(w http.ResponseWriter, r *http.Request) {
	// Decode payload
	var req PostProcessMultipleResponseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	concurrencyLimit := make(chan struct{}, 50) // Limit to 50 concurrent goroutines

	// Validate All Requests
	for _, request := range req.Requests {
		wg.Add(1)
		concurrencyLimit <- struct{}{} // Acquire a slot

		go func() {
			defer wg.Done()
			defer func() { <-concurrencyLimit }() // Release the slot

			// Validate payload
			valid, err := c.Service.Item.ValidateItem(request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// Update item response
			err = c.Service.Item.UpdateItemResponse(request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Post callback
			c.postCallback(request.ItemId)
		}()
	}

	w.WriteHeader(http.StatusOK)
}

func (c *itemController) ReassignItem(w http.ResponseWriter, r *http.Request) {
	// Get user info
	user, err := c.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get query params
	params := mux.Vars(r)
	id := params["itemGuid"]
	approverEmail := params["approver"]
	applicationId := params["ApplicationId"]
	applicationModuleId := params["ApplicationModuleId"]
	approveText := params["ApproveText"]
	rejectText := params["RejectText"]

	err = c.Service.Item.UpdateItemApproverEmail(id, approverEmail, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Post callback
	param := ReassignItemCallback{
		Id:                  id,
		ApproverEmail:       approverEmail,
		Username:            user.Email,
		ApplicationId:       applicationId,
		ApplicationModuleId: applicationModuleId,
		ApproveText:         approveText,
		RejectText:          rejectText,
	}
	go c.postCallbackReassignItem(param)
}

func (c *itemController) postCallback(id string) {
	item, err := c.Service.Item.GetItemById(id)
	if err != nil {
		fmt.Println("Error getting item by id: ", id)
		return
	}

	if item.CallbackUrl == "" {
		fmt.Println("No callback url found")
		return
	} else {
		params := model.ResponseCallback{
			ItemId:       id,
			IsApproved:   item.IsApproved,
			Remarks:      item.ApproverRemarks,
			ResponseDate: item.DateResponded,
			RespondedBy:  item.RespondedBy,
		}

		jsonReq, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Error marshalling response callback: ", err)
			return
		}

		token, err := c.Authenticator.GenerateToken()
		if err != nil {
			fmt.Println("Error generating token: ", err)
			return
		}

		req, err := http.NewRequest("POST", item.CallbackUrl, bytes.NewBuffer(jsonReq))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second * 90,
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error posting callback: ", err)
			return
		}

		isCallbackFailed := res.StatusCode != 200

		err = c.Service.Item.UpdateItemCallback(id, isCallbackFailed)
	}
}

func (c *itemController) postCallbackReassignItem(data ReassignItemCallback) {
	res, err := c.Service.Item.GetItemById(data.Id)
	if err != nil {
		fmt.Println("Error getting item by id: ", data.Id)
		return
	}

	if res.ReassignCallbackUrl != "" {
		jsonReq, err := json.Marshal(data)
		if err != nil {
			return
		}

		token, err := c.Authenticator.GenerateToken()
		if err != nil {
			fmt.Println("Error generating token: ", err)
			return
		}

		req, err := http.NewRequest("POST", res.ReassignCallbackUrl, bytes.NewBuffer(jsonReq))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second * 90,
		}
		_, err = client.Do(req)
		if err != nil {
			fmt.Println("Error posting callback: ", err)
			return
		}

	}
}
