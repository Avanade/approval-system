package approvalRequestApprover

import (
	"main/model"
)

type ApprovalRequestApproverService interface {
	InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error
}
