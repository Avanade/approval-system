package approvalRequestApprover

import (
	"main/model"
)

type ApprovalRequestApproverRepository interface {
	InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error
}
