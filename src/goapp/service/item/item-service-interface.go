package item

import (
	"main/model"
)

type ItemService interface {
	GetFailedCallbacks() ([]string, error)
	GetItemById(id string) (*model.Item, error)
	GetAll(itemOptions model.ItemOptions) (model.Response, error)
	InsertItem(item model.ItemInsertRequest) (string, error)
	ItemIsAuthorized(appId, appModuleId, itemId, approverEmail string) (*model.ItemIsAuthorized, error)
	UpdateItemApproverEmail(itemId, approverEmail, username string) error
	UpdateItemCallback(itemId string, isCallbackFailed bool) error
	UpdateItemDateSent(itemId string) error
	UpdateItemResponse(req model.ProcessResponseRequest) error
	ValidateItem(req model.ProcessResponseRequest) (bool, error)
}
