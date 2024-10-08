package application

import (
	"main/model"
)

type ApplicationService interface {
	GetApplicationById(id string) (*model.Application, error)
}
