package approvalRequestApprover

import (
	"main/model"
)

type ApprovalRequestApproverRepository interface {
	InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error
	GetApproversByItemId(itemId string) ([]string, error)
}
