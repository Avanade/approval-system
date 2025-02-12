package item

import (
	"main/model"
)

type ItemService interface {
	GetAll(itemOptions model.ItemOptions) (model.Response, error)
	GetFailedCallbacks() ([]string, error)
	GetItemById(id string) (*model.Item, error)
	GetItemsForReviewByEmail(email string, page, filter, status int) (*model.Response, error)
	GetByApprover(approver, requestType, organization string, filterOptions model.FilterOptions) (items []model.Item, total int, err error)
	GetInvolvedUsers(itemId string) (*model.InvolvedUsers, error)
	InsertItem(item model.ItemInsertRequest) (string, error)
	ItemIsAuthorized(appId, appModuleId, itemId, approverEmail string) (*model.ItemIsAuthorized, error)
	UpdateItemApproverEmail(itemId, approverEmail, username string) error
	UpdateItemCallback(itemId string, isCallbackFailed bool) error
	UpdateItemDateSent(itemId string) error
	UpdateItemResponse(req model.ProcessResponseRequest) error
	ValidateItem(req model.ProcessResponseRequest) (bool, error)
}
