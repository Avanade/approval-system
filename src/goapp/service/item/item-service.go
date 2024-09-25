package item

import (
	"main/model"
	"main/repository"
)

type itemService struct {
	Repository *repository.Repository
}

func NewItemService(repo *repository.Repository) ItemService {
	return &itemService{
		Repository: repo,
	}
}

func (s *itemService) GetAll(itemOptions model.ItemOptions) (model.Response, error) {
	var result model.Response

	total, err := s.Repository.Item.GetTotalItemsBy(itemOptions)
	if err != nil {
		return model.Response{}, err
	}

	data, err := s.Repository.Item.GetItemsBy(itemOptions)
	if err != nil {
		return model.Response{}, err
	}

	result = model.Response{
		Data:  data,
		Total: total,
	}

	return result, nil
}

func (s *itemService) InsertItem(item model.TypItemInsert) (string, error) {
	id, err := s.Repository.Item.InsertItem(item.ApplicationModuleId, item.Subject, item.Body, item.RequesterEmail)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *itemService) InsertApprovalRequestApprover(approver model.ApprovalRequestApprover) error {
	err := s.Repository.Item.InsertApprovalRequestApprover(approver)
	if err != nil {
		return err
	}
	return nil
}

func (s *itemService) UpdateItemDateSent(itemId string) error {
	err := s.Repository.Item.UpdateItemDateSent(itemId)
	if err != nil {
		return err
	}
	return nil
}
