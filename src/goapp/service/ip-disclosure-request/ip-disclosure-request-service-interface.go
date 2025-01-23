package ipdisclosurerequest

import (
	"main/model"
)

type IpDisclosureRequestService interface {
	InsertIPDisclosureRequest(ipDisclosureRequest *model.IPDRequest) (*model.IPDRequest, error)
	UpdateApprovalRequestId(approvalRequestId string, IPDRequestId int64) error
}
