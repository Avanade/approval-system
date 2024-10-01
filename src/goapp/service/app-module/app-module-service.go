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

func (s *applicationModuleService) GetRequestTypes() ([]RequestType, error) {
	a, err := s.Repository.ApplicationModule.GetAll()
	if err != nil {
		return nil, err
	}

	var result []RequestType
	for _, v := range a {
		result = append(result, RequestType{
			Id:   v.ApplicationModuleId,
			Name: v.ApplicationModuleName,
		})
	}

	return result, nil
}
