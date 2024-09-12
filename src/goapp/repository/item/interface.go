package item

import (
	"main/model"
)

type ItemRepository interface {
	GetItemsBy(itemType model.ItemType, itemStatus model.ItemStatus, requestType, organization, user, search string, offset, filter int) ([]model.Item, error)
	GetTotalItemsBy(itemType model.ItemType, itemStatus model.ItemStatus, requestType, organization, user, search string) (int, error)
}
