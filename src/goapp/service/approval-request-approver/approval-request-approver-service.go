package approvalRequestApprover

import (
	"main/model"
	"main/repository"
)

type approvalRequestApproverService struct {
	Repository *repository.Repository
}

func NewApprovalRequestApproverService(r *repository.Repository) ApprovalRequestApproverService {
	return &approvalRequestApproverService{
		Repository: r,
	}
}

func (s *approvalRequestApproverService) InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error {
	err := s.Repository.ApprovalRequestApprover.InsertApprovalRequestApprover(approver)
	if err != nil {
		return err
	}
	return nil
}
