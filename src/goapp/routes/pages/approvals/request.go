package route

import (
	"encoding/json"
	"fmt"
	"main/model"
	"main/pkg/email"
	"main/pkg/sql"
	"net/http"
	"os"

	rtApi "main/routes/apis"
)

func ApprovalRequestHandler(w http.ResponseWriter, r *http.Request) {

	// Decode payload
	var req model.TypRequestApproval
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// Validate payload
	params := make(map[string]interface{})
	params["ApplicationId"] = req.ApplicationId
	params["ApplicationModuleId"] = req.ApplicationModuleId
	verification, err := db.ExecuteStoredProcedureWithResult("PR_Application_IsValid", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if verification[0]["IsValid"] == "1" {
		for k := range params {
			delete(params, k)
		}
		// Get application module
		params["Id"] = req.ApplicationModuleId
		appModule, err := db.ExecuteStoredProcedureWithResult("PR_ApplicationModules_Select_ById", params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add item to database
		for k := range params {
			delete(params, k)
		}
		params["ApplicationModuleId"] = req.ApplicationModuleId
		params["Subject"] = req.Subject
		params["Body"] = req.Body
		params["RequesterEmail"] = req.RequesterEmail
		item, err := db.ExecuteStoredProcedureWithResult("PR_Items_Insert", params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Compose email message
		emailBodyData := email.TypEmailData{
			Subject:     req.Subject,
			Body:        req.Body,
			ApproveText: fmt.Sprintf("%s", appModule[0]["ApproveText"]),
			RejectText:  fmt.Sprintf("%s", appModule[0]["RejectText"]),
			ApproveUrl:  fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("HOME_URL"), req.ApplicationId, req.ApplicationModuleId, item[0]["Id"]),
			RejectUrl:   fmt.Sprintf("%s/response/%s/%s/%s/0", os.Getenv("HOME_URL"), req.ApplicationId, req.ApplicationModuleId, item[0]["Id"]),
		}

		emailBody, err := email.ComposeEmail(emailBodyData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var toRecipients []email.Recipient

		for _, to := range req.Emails {
			err := rtApi.InsertApprovalRequestApprover(rtApi.ApprovalRequestApprover{
				ItemId:        item[0]["Id"].(string),
				ApproverEmail: to,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			toRecipients = append(toRecipients, email.Recipient{
				Email: to,
			})
		}

		// Send email to approver
		m := email.Message{
			Subject: req.Subject,
			Body: email.Body{
				Content: emailBody,
				Type:    email.HtmlMessageType,
			},
			ToRecipients: toRecipients,
		}

		err = email.SendEmail(m, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update Date Sent
		for k := range params {
			delete(params, k)
		}

		params["Id"] = item[0]["Id"]
		_, err = db.ExecuteStoredProcedure("PR_Items_Update_DateSent", params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["itemId"] = fmt.Sprintf("%v", item[0]["Id"])
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)

	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
