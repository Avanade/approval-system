package appModule

import (
	"main/model"
	"main/repository"
)

type applicationModuleService struct {
	Repository *repository.Repository
}

func NewApplicationModuleService(repo *repository.Repository) ApplicationModuleService {
	return &applicationModuleService{
		Repository: repo,
	}
}

func (s *applicationModuleService) GetApplicationModuleByIdAndApplicationId(applicationId string, applicationModuleId string) (*model.ApplicationModule, error) {
	result, err := s.Repository.ApplicationModule.GetApplicationModuleByIdAndApplicationId(applicationId, applicationModuleId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
