package ipdisclosurerequest

import (
	"main/model"
)

type IpDisclosureRequestService interface {
	GetIPDRequestByApprovalRequestId(approvalRequestId string) (*model.IPDRequest, error)
	InsertIPDisclosureRequest(ipDisclosureRequest *model.IPDRequest) (*model.IPDRequest, error)
	UpdateApprovalRequestId(approvalRequestId string, IPDRequestId int64) error
	UpdateResponse(data *model.ResponseCallback) error
}
