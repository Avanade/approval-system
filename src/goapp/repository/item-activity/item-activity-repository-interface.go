package itemActivity

import (
	"main/model"
)

type ItemActivityRepository interface {
	GetItemActivityByItemId(id string) ([]model.ItemActivity, error)
	AddItemActivity(itemActivity *model.ItemActivity) (*model.ItemActivity, error)
}
