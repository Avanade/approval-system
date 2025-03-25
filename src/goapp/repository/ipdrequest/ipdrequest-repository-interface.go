package ipdrequest

import (
	"main/model"
)

type IpdRequestRepository interface {
	GetIpdRequestByApprovalRequestId(approvalRequestId string) (*model.IPDRequest, error)
	InsertIpdRequest(ipdRequest *model.IPDRequest) (int64, error)
	UpdateApprovalRequestId(approvalRequestId string, IPDRequestId int64) error
	UpdateResponse(data *model.ResponseCallback) error
}
