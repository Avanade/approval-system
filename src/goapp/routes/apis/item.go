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
	Application      string `json:"application"`
	ApproverRemarks  string `json:"approverRemarks"`
	Body             string `json:"body"`
	Created          string `json:"created"`
	DateResponded    string `json:"dateResponded"`
	DateSent         string `json:"dateSent"`
	IsApproved       bool   `json:"isApproved"`
	Module           string `json:"module"`
	Subject          string `json:"subject"`
	ApproveText      string `json:"approveText"`
	RejectText       string `json:"rejectText"`
	ApproveUrl       string `json:"approveUrl"`
	RejectUrl        string `json:"rejectUrl"`
	AllowReassign    bool   `json:"allowReassign"`
	AllowReassignUrl string `json:"allowReassignUrl"`
}

type Response struct {
	Data  []Item `json:"data"`
	Total int    `json:"total"`
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	params := r.URL.Query()

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

	var result Response

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		data, total, err := GetItemsBy(ItemType(itemType), ItemStatus(itemStatus), user, search, offset, filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result = Response{
			Data:  data,
			Total: total,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetItemsBy(itemType ItemType, itemStatus ItemStatus, user, search string, offset, filter int) ([]Item, int, error) {
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING"),
	}

	db, errInit := sql.Init(dbConnectionParam)
	if errInit != nil {
		return []Item{}, 0, errInit
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

	params["Search"] = search

	resultTotal, errResultTotal := db.ExecuteStoredProcedureWithResult("PR_Items_Total", params)
	if errResultTotal != nil {
		return []Item{}, 0, errResultTotal
	}

	params["Offset"] = offset
	params["Filter"] = filter

	total, errTotal := strconv.Atoi(fmt.Sprint(resultTotal[0]["Total"]))
	if errTotal != nil {
		return []Item{}, 0, errTotal
	}

	result, errResult := db.ExecuteStoredProcedureWithResult("PR_Items_Select", params)
	if errResult != nil {
		return []Item{}, 0, errResult
	}

	var items []Item

	for _, v := range result {

		item := Item{
			Application:   v["Application"].(string),
			Created:       v["Created"].(time.Time).String(),
			Module:        v["Module"].(string),
			ApproveText:   v["ApproveText"].(string),
			RejectText:    v["RejectText"].(string),
			AllowReassign: v["AllowReassign"].(bool),
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
			item.AllowReassignUrl = fmt.Sprintf("/responseReassigned/%s/%s/%s/1/%s/%s", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"], v["ApproveText"].(string), v["RejectText"].(string))

		}

		if v["Subject"] != nil {
			item.Subject = v["Subject"].(string)
		}

		items = append(items, item)
	}

	return items, total, nil
}
