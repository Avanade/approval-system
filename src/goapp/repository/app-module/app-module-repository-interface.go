package appModule

import (
	"main/model"
)

type ApplicationModuleRepository interface {
	GetApplicationModuleByIdAndApplicationId(applicationId string, applicationModuleId string) (*model.ApplicationModule, error)
	GetAll() ([]model.ApplicationModule, error)
	IsAuthRequired(applicationModuleId string) (bool, error)
}
