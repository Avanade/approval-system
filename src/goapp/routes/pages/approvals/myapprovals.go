package route

import (
	"fmt"
	"main/pkg/session"
	"main/pkg/sql"

	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func ReAssignApproverHandler(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {

		return
	}
	user := fmt.Sprintf("%s", profile["preferred_username"])
	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	params := mux.Vars(r)

	id := params["itemGuid"]
	approverEmail := params["approver"]

	ApplicationId := params["ApplicationId"]
	ApplicationModuleId := params["ApplicationModuleId"]
	itemId := params["itemId"]
	ApproveText := params["ApproveText"]
	RejectText := params["RejectText"]
	param := map[string]interface{}{

		"Id":            id,
		"ApproverEmail": approverEmail,
		"Username":      user,
	}

	_, err2 := db.ExecuteStoredProcedure("PR_Items_Update_ApproverEmail", param)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	go PostReassignCallback(approverEmail, user, id, ApplicationId, ApplicationModuleId, itemId, ApproveText, RejectText)
}
func itemMapper(item map[string]interface{}, isApproved bool) TypItem {
	var approveUrl string
	var rejectUrl string
	var reassignUrl string
	if !isApproved {
		approveUrl = fmt.Sprintf("/response/%s/%s/%s/1", item["ApplicationId"], item["ApplicationModuleId"], item["ItemId"])
		rejectUrl = fmt.Sprintf("/response/%s/%s/%s/0", item["ApplicationId"], item["ApplicationModuleId"], item["ItemId"])
	}
	reassignUrl = fmt.Sprintf("/responsereassigned/%s/%s/%s/1/%s/%s", item["ApplicationId"], item["ApplicationModuleId"], item["ItemId"], item["ApproveText"], item["RejectText"])

	return TypItem{
		Application:     item["Application"],
		ApproverRemarks: item["ApproverRemarks"],
		Body:            item["Body"],
		Created:         item["Created"],
		DateResponded:   item["DateResponded"],
		DateSent:        item["DateSent"],
		IsApproved:      item["IsApproved"],
		Module:          item["Module"],
		Subject:         item["Subject"],
		ApproveText:     item["ApproveText"],
		RejectText:      item["RejectText"],
		ApproveUrl:      approveUrl,
		RejectUrl:       rejectUrl,
		ReAssignUrl:     reassignUrl,
		Show:            false,
	}
}

type TypHomeData struct {
	Approved            []TypItem
	Pending             []TypItem
	ExportUrl           string
	OrganizationTypeUrl string
}

type TypItem struct {
	Application     interface{}
	ApproverRemarks interface{}
	Body            interface{}
	Created         interface{}
	DateResponded   interface{}
	DateSent        interface{}
	IsApproved      interface{}
	Module          interface{}
	Subject         interface{}
	ApproveText     interface{}
	RejectText      interface{}
	ItemId          interface{}
	ApproveUrl      interface{}
	RejectUrl       interface{}
	ReAssignUrl     interface{}
	Show            bool
}
