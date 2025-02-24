package involvement

import (
	"main/model"
	"main/repository"
)

type involvementService struct {
	Repository *repository.Repository
}

func NewInvolvementService(repo *repository.Repository) InvolvementService {
	return &involvementService{
		Repository: repo,
	}
}

func (s *involvementService) GetInvolvementList() ([]model.Involvement, error) {
	involvementList, err := s.Repository.Involvement.GetInvolvementList()
	if err != nil {
		return []model.Involvement{}, err
	}
	return involvementList, nil
}
