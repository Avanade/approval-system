package application

import (
	"main/model"
	"main/repository"
)

type applicationService struct {
	Repository *repository.Repository
}

func NewApplicationService(r *repository.Repository) ApplicationService {
	return &applicationService{
		Repository: r,
	}
}

func (s *applicationService) GetApplicationById(id string) (*model.Application, error) {
	application, err := s.Repository.Application.GetApplicationById(id)
	if err != nil {
		return nil, err
	}
	return application, nil
}
