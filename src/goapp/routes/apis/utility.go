package routes

import (
	"log"
	"main/pkg/sql"
	"net/http"
	"os"
)

func FillOutApprovalRequestApprovers(w http.ResponseWriter, r *http.Request) {
	items, err := GetAllItems()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, item := range items {
		id := item["Id"].(string)
		err = InsertApprovalRequestApprover(ApprovalRequestApprover{
			ItemId:        id,
			ApproverEmail: item["ApproverEmail"].(string),
		})
		if err != nil {
			log.Println(err.Error())
		}

		if item["DateResponded"] != nil {
			err = UpdateItemById(id, item["ApproverEmail"].(string))
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

type ApprovalRequestApprover struct {
	ItemId        string
	ApproverEmail string
}

func GetAllItems() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_All", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertApprovalRequestApprover(approvalRequestApprover ApprovalRequestApprover) error {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ItemId":        approvalRequestApprover.ItemId,
		"ApproverEmail": approvalRequestApprover.ApproverEmail,
	}

	_, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalRequestApprovers_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func UpdateItemById(id string, respondedBy string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":          id,
		"RespondedBy": respondedBy,
	}

	_, err := db.ExecuteStoredProcedure("PR_Items_UpdateRespondedBy_ById", param)
	if err != nil {
		return err
	}

	return nil
}

func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}
