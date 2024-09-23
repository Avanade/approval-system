package item

import (
	"main/model"
)

type ItemRepository interface {
	GetItemsBy(itemOptions model.ItemOptions) ([]model.Item, error)
	GetTotalItemsBy(itemOptions model.ItemOptions) (int, error)
}
