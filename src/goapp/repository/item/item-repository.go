package item

import (
	"database/sql"
	"fmt"
	"main/model"
	"main/repository"
	"strconv"
	"time"
)

type itemRepository struct {
	repository.Database
}

func NewItemRepository(db repository.Database) ItemRepository {
	return &itemRepository{
		Database: db,
	}
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

	result, err := r.RowsToMap(resList)
	if err != nil {
		return []model.Item{}, err
	}

	var items []model.Item

	for _, v := range result {

		item := model.Item{
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
		} else {
			item.ApproveUrl = fmt.Sprintf("/response/%s/%s/%s/1", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"])
			item.RejectUrl = fmt.Sprintf("/response/%s/%s/%s/0", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"])
			item.AllowReassignUrl = fmt.Sprintf("/responsereassigned/%s/%s/%s/1/%s/%s", v["ApplicationId"], v["ApplicationModuleId"], v["ItemId"], v["ApproveText"].(string), v["RejectText"].(string))

		}

		if v["Subject"] != nil {
			item.Subject = v["Subject"].(string)
		}

		if v["RespondedBy"] != nil {
			item.RespondedBy = v["RespondedBy"].(string)
		}

		rowApprovers, err := r.Query("PR_ApprovalRequestApprovers_Select_ByItemId", sql.Named("ItemId", v["ItemId"].(string)))
		if err != nil {
			return []model.Item{}, err
		}

		approvers, err := r.RowsToMap(rowApprovers)
		if err != nil {
			return []model.Item{}, err
		}

		for _, approver := range approvers {
			item.Approvers = append(item.Approvers, approver["ApproverEmail"].(string))
		}

		items = append(items, item)
	}

	return items, nil
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
