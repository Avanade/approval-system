package involvement

import (
	"main/model"
)

type InvolvementService interface {
	GetInvolvementList() ([]model.Involvement, error)
}
