package item

import (
	"main/model"
)

type ItemService interface {
	GetAll(itemOptions model.ItemOptions) (model.Response, error)
}
