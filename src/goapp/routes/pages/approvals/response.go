package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/pkg/sql"
	"net/http"
	"os"
	"time"
)

func connectSql() (db *sql.DB) {
	db, err := sql.Init(sql.ConnectionParam{ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING")})
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
	return
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
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

	callbackUrl := res[0]["CallbackUrl"].(string)

	if callbackUrl != "" {
		postParams := TypPostParams{
			ItemId:       itemId,
			IsApproved:   res[0]["IsApproved"].(bool),
			Remarks:      res[0]["ApproverRemarks"].(string),
			ResponseDate: approvalDate.Format("2006-01-02T15:04:05.000Z"),
			RespondedBy:  res[0]["RespondedBy"].(string),
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
	if err != nil {
		ch <- nil
	}

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
	RespondedBy  string `json:"respondedBy"`
}
