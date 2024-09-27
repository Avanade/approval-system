package appModule

import (
	"main/model"
)

type ApplicationModuleRepository interface {
	GetApplicationModuleByIdAndApplicationId(applicationId string, applicationModuleId string) (*model.ApplicationModule, error)
}
