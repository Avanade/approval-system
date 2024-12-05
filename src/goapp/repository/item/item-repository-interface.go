package item

import (
	"main/model"
)

type ItemRepository interface {
	GetFailedCallbacks() ([]string, error)
	GetItemById(id string) (*model.Item, error)
	GetItemsByApprover(approver, requestType, organization string, filterOptions model.FilterOptions) (items []model.Item, total int, err error)
	GetItemsBy(itemOptions model.ItemOptions) ([]model.Item, error)
	GetTotalItemsBy(itemOptions model.ItemOptions) (int, error)
	InsertItem(appModuleId, subject, body, requesterEmail string) (string, error)
	ItemIsAuthorized(appId, appModuleId, itemId, approverEmail string) (*model.ItemIsAuthorized, error)
	UpdateItemApproverEmail(id, approverEmail, username string) error
	UpdateItemCallback(id string, isCallbackFailed bool) error
	UpdateItemDateSent(id string) error
	UpdateItemResponse(id, remarks, email string, isApproved bool) error
	ValidateItem(appId, appModuleId, itemId, email string) (bool, error)
}
