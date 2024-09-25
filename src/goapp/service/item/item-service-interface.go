package item

import (
	"main/model"
)

type ItemService interface {
	GetAll(itemOptions model.ItemOptions) (model.Response, error)
	InsertItem(item model.TypItemInsert) (string, error)
	InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error
	UpdateItemDateSent(itemId string) error
}
