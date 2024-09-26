package approvalRequestApprover

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type approvalRequestApproverRepository struct {
	db.Database
}

func NewApprovalRequestApproverRepository(db db.Database) ApprovalRequestApproverRepository {
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
