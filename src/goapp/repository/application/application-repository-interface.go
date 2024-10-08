package application

import (
	"main/model"
)

type ApplicationRepository interface {
	GetApplicationById(id string) (*model.Application, error)
}
