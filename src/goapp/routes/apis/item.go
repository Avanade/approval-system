package routes

import (
	"encoding/json"
	"fmt"
	"main/pkg/sql"
	"net/http"
	"os"
)

type ItemType int8

const (
	RequestItem ItemType = iota
	ApprovalItem
	AllType
)

type ItemStatus int8

const (
	Pending ItemStatus = iota
	Approved
	Rejected
	Closed // Disapproved, Approved
	All    // Disapproved, Approved, Pending
)

type Item struct {
	Application      string   `json:"application"`
	ApproverRemarks  string   `json:"approverRemarks"`
	Body             string   `json:"body"`
	Created          string   `json:"created"`
	DateResponded    string   `json:"dateResponded"`
	DateSent         string   `json:"dateSent"`
	IsApproved       bool     `json:"isApproved"`
	Module           string   `json:"module"`
	Subject          string   `json:"subject"`
	ApproveText      string   `json:"approveText"`
	RejectText       string   `json:"rejectText"`
	ApproveUrl       string   `json:"approveUrl"`
	RejectUrl        string   `json:"rejectUrl"`
	AllowReassign    bool     `json:"allowReassign"`
	AllowReassignUrl string   `json:"allowReassignUrl"`
	RespondedBy      string   `json:"respondedBy"`
	Approvers        []string `json:"approvers"`
	RequestedBy      string   `json:"requestedBy"`
}

type Response struct {
	Data  []Item `json:"data"`
	Total int    `json:"total"`
}

type RequestType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetRequestTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, err := GetRequestType()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetRequestType() ([]RequestType, error) {
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApplicationModules_Select", nil)
	if err != nil {
		return nil, err
	}

	var requestTypes []RequestType

	for _, v := range result {
		requestType := RequestType{
			Id:   fmt.Sprintf("%v", v["Id"]),
			Name: fmt.Sprintf("%v", v["Name"]),
		}
		requestTypes = append(requestTypes, requestType)
	}

	return requestTypes, nil
}
