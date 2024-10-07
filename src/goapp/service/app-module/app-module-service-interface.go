package appModule

import (
	"main/model"
)

type ApplicationModuleService interface {
	GetApplicationModuleByIdAndApplicationId(applicationId string, applicationModuleId string) (*model.ApplicationModule, error)
	GetRequestTypes() ([]RequestType, error)
}
