package routes

import (
	"encoding/json"
	"fmt"
	"main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ItemType int8

const (
	RequestItem ItemType = iota
	ApprovalItem
	AllType
)

type ItemStatus int8

const (
	Disapproved ItemStatus = iota
	Approved
	Pending
	AllStatus
)

type Item struct {
	Application     string `json:"application"`
	ApproverRemarks string `json:"approverRemarks"`
	Body            string `json:"body"`
	Created         string `json:"created"`
	DateResponded   string `json:"dateResponded"`
	DateSent        string `json:"dateSent"`
	IsApproved      bool   `json:"isApproved"`
	Module          string `json:"module"`
	Subject         string `json:"subject"`
	ApproveText     string `json:"approveText"`
	RejectText      string `json:"rejectText"`
	ApproveUrl      string `json:"approverUrl"`
	RejectUrl       string `json:"rejectUrl"`
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	session, errAuth := session.Store.Get(r, "auth-session")
	if errAuth != nil {
		http.Error(w, errAuth.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, errAuth.Error(), http.StatusInternalServerError)
		return
	}
	user := fmt.Sprintf("%s", profile["preferred_username"])

	itemType, errItemType := strconv.ParseInt(vars["type"], 10, 64)
	if errItemType != nil {
		http.Error(w, errItemType.Error(), http.StatusInternalServerError)
		return
	}

	itemStatus, errItemStatus := strconv.ParseInt(vars["status"], 10, 64)
	if errItemStatus != nil {
		http.Error(w, errItemStatus.Error(), http.StatusInternalServerError)
		return
	}

	result, err := GetItemsBy(user, ItemType(itemType), ItemStatus(itemStatus))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetItemsBy(user string, itemType ItemType, itemStatus ItemStatus) ([]Item, error) {
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING"),
	}

	db, errInit := sql.Init(dbConnectionParam)
	if errInit != nil {
		return []Item{}, errInit
	}
	defer db.Close()

	params := make(map[string]interface{})
	if itemType != AllType {
		params["ItemType"] = itemType
		params["User"] = user
	}
	if itemStatus != AllStatus {
		if itemStatus == Pending {
			params["IsApproved"] = nil
		} else {
			params["IsApproved"] = itemStatus
		}
	}
	result, errExec := db.ExecuteStoredProcedureWithResult("PR_Items_Select", params)
	if errExec != nil {
		return []Item{}, errExec
	}

	var items []Item

	for _, v := range result {
		item := Item{
			Application: v["Application"].(string),
			Created:     v["Created"].(time.Time).String(),
			Module:      v["Module"].(string),
			ApproveText: v["ApproveText"].(string),
			RejectText:  v["RejectText"].(string),
		}

		if v["ApproverRemarks"] != nil {
			item.ApproverRemarks = v["ApproverRemarks"].(string)
		}

		if v["Body"] != nil {
			item.Body = v["Body"].(string)
		}

		if v["DateResponded"] != nil {
			item.DateResponded = v["DateResponded"].(time.Time).String()
		}

		if v["DateSent"] != nil {
			item.DateSent = v["DateSent"].(time.Time).String()
		}

		if v["IsApproved"] != nil {
			item.IsApproved = v["IsApproved"].(bool)
		} else {
			item.ApproveUrl = fmt.Sprintf("/response/%s/%s/%s/1", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"])
			item.RejectUrl = fmt.Sprintf("/response/%s/%s/%s/0", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"])
		}

		if v["Subject"] != nil {
			item.Subject = v["Subject"].(string)
		}

		items = append(items, item)
	}

	return items, nil
}
