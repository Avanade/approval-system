package involvement

import (
	"main/model"
)

type InvolvementRepository interface {
	GetInvolvementList() ([]model.Involvement, error)
}
