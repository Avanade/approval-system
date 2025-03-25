package ipdrInvolvement

import (
	"main/model"
)

type IpdrInvolvementRepository interface {
	GetIpdrInvolvementByRequestId(requestId int64) ([]string, []string, error)
	InsertIpdrInvolvement(ipdrInvolvement model.IpdrInvolvement) error
}
