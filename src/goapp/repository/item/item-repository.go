package item

import (
	"database/sql"
	"fmt"
	db "main/infrastructure/database"
	"main/model"
	"strconv"
	"strings"
	"time"
)

type itemRepository struct {
	*db.Database
}

func NewItemRepository(db *db.Database) ItemRepository {
	return &itemRepository{
		Database: db,
	}
}

func (r *itemRepository) GetFailedCallbacks() ([]string, error) {
	row, err := r.Query("PR_Items_Select_FailedCallbacks")
	if err != nil {
		return []string{}, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return []string{}, err
	}

	var failedCallbacks []string
	for _, v := range result {
		failedCallbacks = append(failedCallbacks, v["Id"].(string))
	}

	return failedCallbacks, nil
}

func (r *itemRepository) GetItemById(id string) (*model.Item, error) {
	row, err := r.Query("PR_Items_Select_ById", sql.Named("Id", id))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	item := model.Item{
		Id:            id,
		ApplicationId: result[0]["ApplicationId"].(string),
		Application:   result[0]["Application"].(string),
		Module:        result[0]["Module"].(string),
		ApproveText:   result[0]["ApproveText"].(string),
		RejectText:    result[0]["RejectText"].(string),
		Created:       result[0]["Created"].(time.Time).String(),
		RequestedBy:   result[0]["RequestedBy"].(string),
	}

	if result[0]["ApproverRemarks"] != nil {
		item.ApproverRemarks = result[0]["ApproverRemarks"].(string)
	}

	if result[0]["Body"] != nil {
		item.Body = result[0]["Body"].(string)
	}

	if result[0]["DateResponded"] != nil {
		item.DateResponded = result[0]["DateResponded"].(time.Time).String()
	}

	if result[0]["DateSent"] != nil {
		item.DateSent = result[0]["DateSent"].(time.Time).String()
	}

	if result[0]["IsApproved"] != nil {
		item.IsApproved = result[0]["IsApproved"].(bool)
	}

	if result[0]["Subject"] != nil {
		item.Subject = result[0]["Subject"].(string)
	}

	if result[0]["CallbackUrl"] != nil {
		item.CallbackUrl = result[0]["CallbackUrl"].(string)
	}

	if result[0]["ReassignCallbackUrl"] != nil {
		item.ReassignCallbackUrl = result[0]["ReassignCallbackUrl"].(string)
	}

	if result[0]["RespondedBy"] != nil {
		item.RespondedBy = result[0]["RespondedBy"].(string)
	}

	if result[0]["ApplicationModuleId"] != nil {
		item.ModuleId = result[0]["ApplicationModuleId"].(string)
	}

	return &item, nil
}

func (r *itemRepository) GetItemsByApprover(approver, requestType, organization string, filterOptions model.FilterOptions) (items []model.Item, total int, err error) {
	var params []interface{}

	if requestType != "" {
		params = append(params, sql.Named("RequestType", requestType))
	}

	if organization != "" {
		params = append(params, sql.Named("Organization", organization))
	}

	if filterOptions.Filter != 0 {
		params = append(params, sql.Named("Filter", filterOptions.Filter))
	}

	if filterOptions.Page != 0 {
		offset := filterOptions.Page * filterOptions.Filter
		params = append(params, sql.Named("Offset", offset))
	}

	params = append(params, sql.Named("Approver", approver))

	rows, err := r.Query("PR_Items_Select_ByApprover", params...)
	if err != nil {
		return []model.Item{}, 0, err
	}
	defer rows.Close()

	result, err := r.RowsToMap(rows)
	if err != nil {
		return []model.Item{}, 0, err
	}

	for _, v := range result {
		item := model.Item{
			Id:            v["Id"].(string),
			Subject:       v["Subject"].(string),
			Application:   v["ApplicationName"].(string),
			ApplicationId: v["ApplicationId"].(string),
			Module:        v["ApplicationModuleName"].(string),
			ModuleId:      v["ApplicationModuleId"].(string),
			Created:       v["Created"].(time.Time).String(),
			RequestedBy:   v["RequestedBy"].(string),
			Body:          v["Body"].(string),
		}
		approvers := v["Approvers"].(string)
		approversArray := strings.Split(approvers, ",")
		if len(approversArray) > 0 {
			item.Approvers = approversArray
		}

		items = append(items, item)
	}

	if rows.NextResultSet() {
		if rows.Next() {
			err = rows.Scan(&total)
			if err != nil {
				return []model.Item{}, 0, err
			}
		}
	}

	return items, total, nil
}

func (r *itemRepository) GetItemsByModuleId(moduleId string, filterOptions model.FilterOptions, status int) ([]model.Item, error) {
	var items []model.Item
	offset := filterOptions.Page
	if filterOptions.Page != 0 {
		offset = filterOptions.Page * filterOptions.Filter
	}

	row, err := r.Query("PR_Items_Select_ByModuleId",
		sql.Named("ModuleId", moduleId),
		sql.Named("Offset", offset),
		sql.Named("Filter", filterOptions.Filter),
		sql.Named("Search", filterOptions.Search),
		sql.Named("IsApproved", status))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		item := model.Item{
			Id:          v["ItemId"].(string),
			Application: v["Application"].(string),
			Module:      v["Module"].(string),
			Created:     v["Created"].(time.Time).String(),
			RequestedBy: v["RequestedBy"].(string),
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
		}

		if v["Subject"] != nil {
			item.Subject = v["Subject"].(string)
		}

		if v["RespondedBy"] != nil {
			item.RespondedBy = v["RespondedBy"].(string)
		}

		if v["ApplicationModuleId"] != nil {
			item.ModuleId = v["ApplicationModuleId"].(string)
		}

		if v["ApplicationId"] != nil {
			item.ApplicationId = v["ApplicationId"].(string)
		}

		if v["IPTitle"] != nil {
			item.IPDRTitle = v["IPTitle"].(string)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *itemRepository) GetItemsBy(itemOptions model.ItemOptions) ([]model.Item, error) {
	var params []interface{}

	if model.ItemType(itemOptions.ItemType) != model.AllType {
		params = append(params, sql.Named("ItemType", itemOptions.ItemType))
		params = append(params, sql.Named("User", itemOptions.User))
	}

	if itemOptions.RequestType != "" {
		params = append(params, sql.Named("RequestType", itemOptions.RequestType))
	}

	if itemOptions.Organization != "" {
		params = append(params, sql.Named("Organization", itemOptions.Organization))
	}

	params = append(params, sql.Named("IsApproved", itemOptions.ItemStatus))
	params = append(params, sql.Named("Search", itemOptions.Search))
	params = append(params, sql.Named("Offset", itemOptions.Offset))
	params = append(params, sql.Named("Filter", itemOptions.Filter))

	resList, err := r.Query("PR_Items_Select", params...)
	if err != nil {
		return []model.Item{}, err
	}
	defer resList.Close()

	result, err := r.RowsToMap(resList)
	if err != nil {
		return []model.Item{}, err
	}

	var items []model.Item

	for _, v := range result {

		item := model.Item{
			Id:            v["ItemId"].(string),
			Application:   v["Application"].(string),
			Created:       v["Created"].(time.Time).String(),
			Module:        v["Module"].(string),
			ApproveText:   v["ApproveText"].(string),
			RejectText:    v["RejectText"].(string),
			AllowReassign: v["AllowReassign"].(bool),
			RequestedBy:   v["RequestedBy"].(string),
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
		}
		item.ApproveUrl = fmt.Sprintf("/response/%s/%s/%s/1", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"])
		item.RejectUrl = fmt.Sprintf("/response/%s/%s/%s/0", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"])
		item.AllowReassignUrl = fmt.Sprintf("/responsereassigned/%s/%s/%s/1/%s/%s", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"], v["ApproveText"].(string), v["RejectText"].(string))

		if v["Subject"] != nil {
			item.Subject = v["Subject"].(string)
		}

		if v["RespondedBy"] != nil {
			item.RespondedBy = v["RespondedBy"].(string)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *itemRepository) GetTotalItemsByModuleId(appModuleId string, status int) (int, error) {
	row, err := r.QueryRow("PR_Items_Total_ByModuleId",
		sql.Named("IsApproved", status),
		sql.Named("ModuleId", appModuleId),
	)

	if err != nil {
		return 0, err
	}

	var total int
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *itemRepository) GetTotalItemsBy(itemOptions model.ItemOptions) (int, error) {
	var params []interface{}

	if model.ItemType(itemOptions.ItemType) != model.AllType {
		params = append(params, sql.Named("ItemType", itemOptions.ItemType))
		params = append(params, sql.Named("User", itemOptions.User))
	}

	if itemOptions.RequestType != "" {
		params = append(params, sql.Named("RequestType", itemOptions.RequestType))
	}

	if itemOptions.Organization != "" {
		params = append(params, sql.Named("Organization", itemOptions.Organization))
	}

	params = append(params, sql.Named("IsApproved", itemOptions.ItemStatus))
	params = append(params, sql.Named("Search", itemOptions.Search))

	rowTotal, err := r.Query("PR_Items_Total", params...)
	if err != nil {
		return 0, err
	}
	defer rowTotal.Close()

	resultTotal, err := r.RowsToMap(rowTotal)
	if err != nil {
		return 0, err
	}

	total, err := strconv.Atoi(fmt.Sprint(resultTotal[0]["Total"]))
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *itemRepository) InsertItem(appModuleId, subject, body, requesterEmail string) (string, error) {
	rowItem, err := r.Query("PR_Items_Insert",
		sql.Named("ApplicationModuleId", appModuleId),
		sql.Named("Subject", subject),
		sql.Named("Body", body),
		sql.Named("RequesterEmail", requesterEmail),
	)

	if err != nil {
		return "", err
	}
	defer rowItem.Close()

	resultItem, err := r.RowsToMap(rowItem)
	if err != nil {
		return "", err
	}

	return resultItem[0]["Id"].(string), nil
}

func (r *itemRepository) UpdateItemApproverEmail(id, approverEmail, username string) error {
	row, err := r.Query("PR_Items_Update_ApproverEmail",
		sql.Named("Id", id),
		sql.Named("ApproverEmail", approverEmail),
		sql.Named("Username", username),
	)
	if err != nil {
		return err
	}
	defer row.Close()
	return nil
}
func (r *itemRepository) ItemIsAuthorized(appId, appModuleId, itemId, approverEmail string) (*model.ItemIsAuthorized, error) {
	row, err := r.Query("PR_Items_IsAuthorized",
		sql.Named("ApplicationId", appId),
		sql.Named("ApplicationModuleId", appModuleId),
		sql.Named("ItemId", itemId),
		sql.Named("ApproverEmail", approverEmail),
	)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	i := model.ItemIsAuthorized{
		IsAuthorized: result[0]["IsAuthorized"] == "1",
	}

	if result[0]["IsApproved"] != nil {
		i.IsApproved = &model.NullBool{Value: result[0]["IsApproved"].(bool)}
	} else {
		i.IsApproved = nil
	}

	if result[0]["RequireRemarks"] != nil {
		i.RequireRemarks = result[0]["RequireRemarks"].(bool)
	}

	return &i, nil
}

func (r *itemRepository) UpdateItemCallback(id string, isCallbackFailed bool) error {
	row, err := r.Query("PR_Items_Update_Callback",
		sql.Named("ItemId", id),
		sql.Named("IsCallbackFailed", isCallbackFailed),
	)
	if err != nil {
		return err
	}
	defer row.Close()
	return nil
}

func (r *itemRepository) UpdateItemDateSent(id string) error {
	row, err := r.Query("PR_Items_Update_DateSent",
		sql.Named("Id", id),
	)

	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}

func (r *itemRepository) UpdateItemResponse(id, remarks, email string, isApproved bool) error {
	row, err := r.Query("PR_Items_Update_Response",
		sql.Named("Id", id),
		sql.Named("ApproverRemarks", remarks),
		sql.Named("Username", email),
		sql.Named("IsApproved", isApproved),
	)
	if err != nil {
		return err
	}
	defer row.Close()
	return nil
}

func (r *itemRepository) ValidateItem(appId, appModuleId, itemId, email string) (bool, error) {
	row, err := r.Query("PR_Items_IsValid",
		sql.Named("ApplicationId", appId),
		sql.Named("ApplicationModuleId", appModuleId),
		sql.Named("ItemId", itemId),
		sql.Named("ApproverEmail", email),
	)
	if err != nil {
		return false, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return false, err
	}

	return result[0]["IsValid"] == "1", nil
}
