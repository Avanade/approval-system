package ipdrequest

import (
	"main/model"
)

type IpdRequestRepository interface {
	InsertIpdRequest(ipdRequest *model.IPDRequest) (int64, error)
	UpdateApprovalRequestId(approvalRequestId string, IPDRequestId int64) error
}
