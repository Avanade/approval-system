package item

import (
	"main/model"
)

type ItemRepository interface {
	GetItemById(id string) (*model.Item, error)
	GetItemsBy(itemOptions model.ItemOptions) ([]model.Item, error)
	GetTotalItemsBy(itemOptions model.ItemOptions) (int, error)
	InsertItem(appModuleId, subject, body, requesterEmail string) (string, error)
	UpdateItemApproverEmail(id, approverEmail, username string) error
	UpdateItemCallback(id string, isCallbackFailed bool) error
	UpdateItemDateSent(id string) error
	UpdateItemResponse(id, remarks, email string, isApproved bool) error
	ValidateItem(appId, appModuleId, itemId, email string) (bool, error)
}
