package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"
	session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func connectSql() (db *sql.DB) {
	db, err := sql.Init(sql.ConnectionParam{ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING")})
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
	return
}

func handleErrorReturn(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
}
func ResponseReassignedeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var username string
		sessionaz, _ := session.Store.Get(r, "auth-session")
		iprofile := sessionaz.Values["profile"]

		if iprofile != nil {
			profile := iprofile.(map[string]interface{})
			username = profile["preferred_username"].(string)
		}
		params := mux.Vars(r)

		appGuid := params["appGuid"]
		appModuleGuid := params["appModuleGuid"]
		itemGuid := params["itemGuid"]
		isApproved := params["isApproved"]

		sqlParamsIsAuth := map[string]interface{}{
			"ApplicationId":       appGuid,
			"ApplicationModuleId": appModuleGuid,
			"ItemId":              itemGuid,
			"ApproverEmail":       username,
		}

		sqlParamsItems := map[string]interface{}{
			"Id": itemGuid,
		}

		db := connectSql()
		defer db.Close()
		resIsAuth, err := db.ExecuteStoredProcedureWithResult("PR_RESPONSE_IsAuthorized", sqlParamsIsAuth)
		handleErrorReturn(w, err)

		isAuth := resIsAuth[0]["IsAuthorized"]
		if isAuth == "0" {
			template.UseTemplate(&w, r, "Unauthorized", nil)
		} else {
			isProcessed := resIsAuth[0]["IsApproved"]
			if isProcessed != nil {
				var text string
				if isProcessed == true {
					text = "approved"
				} else {
					text = "rejected"
				}
				data := map[string]interface{}{
					"response": text,
				}
				template.UseTemplate(&w, r, "AlreadyProcessed", data)
			} else {
				resItems, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ById", sqlParamsItems)

				handleErrorReturn(w, err)
				requireRemarks := resIsAuth[0]["RequireRemarks"]
				data := map[string]interface{}{
					"ApplicationId":       appGuid,
					"ApplicationModuleId": appModuleGuid,
					"ItemId":              itemGuid,
					"ApproverEmail":       username,
					"IsApproved":          isApproved,
					"Data":                resItems[0],
					"RequireRemarks":      requireRemarks,
				}
				template.UseTemplate(&w, r, "responsereassign", data)
			}

		}
	}
	//template.UseTemplate(&w, r, "responsereassign", nil)
}
func ResponseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var username string
		sessionaz, _ := session.Store.Get(r, "auth-session")
		iprofile := sessionaz.Values["profile"]

		if iprofile != nil {
			profile := iprofile.(map[string]interface{})
			username = profile["preferred_username"].(string)
		}
		params := mux.Vars(r)

		appGuid := params["appGuid"]
		appModuleGuid := params["appModuleGuid"]
		itemGuid := params["itemGuid"]
		isApproved := params["isApproved"]

		sqlParamsIsAuth := map[string]interface{}{
			"ApplicationId":       appGuid,
			"ApplicationModuleId": appModuleGuid,
			"ItemId":              itemGuid,
			"ApproverEmail":       username,
		}

		sqlParamsItems := map[string]interface{}{
			"Id": itemGuid,
		}

		db := connectSql()
		defer db.Close()
		resIsAuth, err := db.ExecuteStoredProcedureWithResult("PR_RESPONSE_IsAuthorized", sqlParamsIsAuth)
		handleErrorReturn(w, err)

		isAuth := resIsAuth[0]["IsAuthorized"]
		if isAuth == "0" {
			template.UseTemplate(&w, r, "Unauthorized", nil)
		} else {
			isProcessed := resIsAuth[0]["IsApproved"]
			if isProcessed != nil {
				var text string
				if isProcessed == true {
					text = "approved"
				} else {
					text = "rejected"
				}
				data := map[string]interface{}{
					"response": text,
				}
				template.UseTemplate(&w, r, "AlreadyProcessed", data)
			} else {
				resItems, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ById", sqlParamsItems)

				handleErrorReturn(w, err)
				requireRemarks := resIsAuth[0]["RequireRemarks"]
				data := map[string]interface{}{
					"ApplicationId":       appGuid,
					"ApplicationModuleId": appModuleGuid,
					"ItemId":              itemGuid,
					"ApproverEmail":       username,
					"IsApproved":          isApproved,
					"Data":                resItems[0],
					"RequireRemarks":      requireRemarks,
				}
				template.UseTemplate(&w, r, "response", data)
			}

		}
	}
}

func ProcessResponseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Decode payload
		var req models.TypRequestProcess
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {

			return
		}

		db := connectSql()
		defer db.Close()

		// Validate payload
		params := make(map[string]interface{})
		params["ApplicationId"] = req.ApplicationId
		params["ApplicationModuleId"] = req.ApplicationModuleId
		params["ItemId"] = req.ItemId
		params["ApproverEmail"] = req.ApproverEmail
		verification, err := db.ExecuteStoredProcedureWithResult("PR_Items_IsValid", params)
		handleErrorReturn(w, err)
		authReq := true
		authReq = session.IsAuthRequired(req.ApplicationModuleId)
		if authReq == false {
			verification[0]["IsValid"] = "1"
		}

		if verification[0]["IsValid"] == "1" {
			for k := range params {
				delete(params, k)
			}
			isApproved, _ := strconv.ParseBool(req.IsApproved)
			params["Id"] = req.ItemId
			params["IsApproved"] = isApproved
			params["ApproverRemarks"] = req.Remarks
			params["Username"] = req.ApproverEmail
			_, err := db.ExecuteStoredProcedure("PR_Items_Update_Response", params)
			handleErrorReturn(w, err)
			postCallback(req.ItemId)
			return
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}

func ProcessFailedCallbacks() {
	db := connectSql()
	defer db.Close()
	res, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_FailedCallbacks", nil)
	handleError(err)

	for _, i := range res {
		go postCallback(i["Id"].(string))
	}
}

func postCallback(itemId string) {
	db := connectSql()
	defer db.Close()

	queryParams := map[string]interface{}{
		"Id": itemId,
	}
	res, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ById", queryParams)
	handleError(err)
	approvalDate := res[0]["DateResponded"].(time.Time)

	var callbackUrl string

	callbackUrl = res[0]["CallbackUrl"].(string)

	if callbackUrl != "" {
		postParams := TypPostParams{
			ItemId:       itemId,
			IsApproved:   res[0]["IsApproved"].(bool),
			Remarks:      res[0]["ApproverRemarks"].(string),
			ResponseDate: approvalDate.Format("2006-01-02T15:04:05.000Z"),
		}

		ch := make(chan *http.Response)

		// var res *http.Response

		go getHttpPostResponseStatus(callbackUrl, postParams, ch)

		res := <-ch

		isCallbackFailed := true
		if res != nil {
			if res.StatusCode == 200 {
				isCallbackFailed = false
			}
		}

		params := map[string]interface{}{
			"ItemId":           itemId,
			"IsCallbackFailed": isCallbackFailed,
		}
		db.ExecuteStoredProcedure("PR_Items_Update_Callback", params)

	}

}

func getHttpPostResponseStatus(url string, data interface{}, ch chan *http.Response) {
	jsonReq, err := json.Marshal(data)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		ch <- nil
	}
	ch <- res
}

type TypPostParams struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"remarks"`
	ResponseDate string `json:"responseDate"`
}
type TypReeassignParams struct {
	Id            string `json:"Id"`
	ApproverEmail string `json:"ApproverEmail"`
	Username      string `json:"Username"`
}

func PostReassignCallback(userEmail string, user string, itemId string) {

	db := connectSql()
	defer db.Close()

	queryParams := map[string]interface{}{
		"Id": itemId,
	}
	res, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ById", queryParams)
	handleError(err)

	var ReassignCallbackUrl string

	ReassignCallbackUrl = res[0]["ReassignCallbackUrl"].(string)
	if ReassignCallbackUrl != "" {
		postParams := TypReeassignParams{
			Id:            itemId,
			ApproverEmail: userEmail,
			Username:      user,
		}

		ch2 := make(chan *http.Response)
		// var res *http.Response

		getHttpPostResponseStatus(ReassignCallbackUrl, postParams, ch2)
		res := <-ch2

		isCallbackFailed := true
		if res != nil {
			if res.StatusCode == 200 {
				isCallbackFailed = false
			}
		}
		fmt.Println(isCallbackFailed)
		// res2 := <-ch2
		// isCallbackFailed := true
		// if res2 != nil {
		// 	if res2.StatusCode == 200 {
		// 		isCallbackFailed = false
		// 	}
		// }
		// params := map[string]interface{}{
		// 	"ItemId":           itemId,
		// 	"IsCallbackFailed": isCallbackFailed,
		// }
		// db.ExecuteStoredProcedure("PR_Items_Update_Callback", params)

	}

}
