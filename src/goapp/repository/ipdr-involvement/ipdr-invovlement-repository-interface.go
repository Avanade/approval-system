package ipdrInvolvement

import (
	"main/model"
)

type IpdrInvolvementRepository interface {
	InsertIpdrInvolvement(ipdrInvolvement model.IpdrInvolvement) error
}