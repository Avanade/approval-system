package approvalRequestApprover

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type approvalRequestApproverRepository struct {
	*db.Database
}

func NewApprovalRequestApproverRepository(db *db.Database) ApprovalRequestApproverRepository {
	return &approvalRequestApproverRepository{
		Database: db,
	}
}

func (r *approvalRequestApproverRepository) InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error {
	_, err := r.Query("PR_ApprovalRequestApprovers_Insert",
		sql.Named("ItemId", approver.ItemId),
		sql.Named("ApproverEmail", approver.ApproverEmail),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *approvalRequestApproverRepository) GetApproversByItemId(itemId string) ([]string, error) {
	var result []string
	rowApprovers, err := r.Query("PR_ApprovalRequestApprovers_Select_ByItemId", sql.Named("ItemId", itemId))
	if err != nil {
		return nil, err
	}

	approvers, err := r.RowsToMap(rowApprovers)
	if err != nil {
		return nil, err
	}

	for _, approver := range approvers {
		result = append(result, approver["ApproverEmail"].(string))
	}

	return result, nil
}
