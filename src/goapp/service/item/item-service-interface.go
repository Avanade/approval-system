package item

import (
	"main/model"
)

type ItemService interface {
	GetItemById(id string) (*model.Item, error)
	GetAll(itemOptions model.ItemOptions) (model.Response, error)
	InsertItem(item model.ItemInsertRequest) (string, error)
	UpdateItemApproverEmail(itemId, approverEmail, username string) error
	UpdateItemCallback(itemId string, isCallbackFailed bool) error
	UpdateItemDateSent(itemId string) error
	UpdateItemResponse(req model.ProcessResponseRequest) error
	ValidateItem(req model.ProcessResponseRequest) (bool, error)
}
