package ipdisclosurerequest

import (
	"main/model"
	"main/repository"
)

type ipDisclosureRequestService struct {
	Repository *repository.Repository
}

func NewIpDisclosureRequestService(repo *repository.Repository) IpDisclosureRequestService {
	return &ipDisclosureRequestService{
		Repository: repo,
	}
}

func (s *ipDisclosureRequestService) GetIPDRequestByApprovalRequestId(approvalRequestId string) (*model.IPDRequest, error) {
	item, err := s.Repository.IPDRequest.GetIpdRequestByApprovalRequestId(approvalRequestId)
	if err != nil {
		return nil, err
	}

	involvementIds, involvements, err := s.Repository.IpdrInvolvement.GetIpdrInvolvementByRequestId(item.RequestId)
	if err != nil {
		return nil, err
	}

	item.InvolvementId = involvementIds
	item.Involvement = involvements

	return item, nil
}

func (s *ipDisclosureRequestService) InsertIPDisclosureRequest(ipDisclosureRequest *model.IPDRequest) (*model.IPDRequest, error) {
	id, err := s.Repository.IPDRequest.InsertIpdRequest(ipDisclosureRequest)
	if err != nil {
		return nil, err
	}

	for _, involvementId := range ipDisclosureRequest.InvolvementId {
		data := model.IpdrInvolvement{
			RequestId:     id,
			InvolvementId: involvementId,
		}
		err = s.Repository.IpdrInvolvement.InsertIpdrInvolvement(data)
		if err != nil {
			return nil, err
		}
	}

	ipDisclosureRequest.RequestId = id

	return ipDisclosureRequest, nil
}

func (s *ipDisclosureRequestService) UpdateApprovalRequestId(approvalRequestId string, IPDRequestId int64) error {
	return s.Repository.IPDRequest.UpdateApprovalRequestId(approvalRequestId, IPDRequestId)
}

func (s *ipDisclosureRequestService) UpdateResponse(data *model.ResponseCallback) error {
	return s.Repository.IPDRequest.UpdateResponse(data)
}
