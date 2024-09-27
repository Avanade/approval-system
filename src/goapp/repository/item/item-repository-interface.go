package item

import (
	"main/model"
)

type ItemRepository interface {
	GetItemsBy(itemOptions model.ItemOptions) ([]model.Item, error)
	GetTotalItemsBy(itemOptions model.ItemOptions) (int, error)
	InsertItem(appModuleId, subject, body, requesterEmail string) (string, error)
	UpdateItemDateSent(id string) error
}
