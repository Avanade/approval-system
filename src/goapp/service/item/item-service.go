package item

import (
	"main/config"
	"main/model"
	"main/repository"
)

type itemService struct {
	Repository    *repository.Repository
	configManager config.ConfigManager
}

func NewItemService(repo *repository.Repository, configManager config.ConfigManager) ItemService {
	return &itemService{
		Repository:    repo,
		configManager: configManager,
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

	for i := 0; i < len(data); i++ {
		if data[i].Approvers == nil {
			continue
		}
		if len(data[i].Approvers) == 0 {
			continue
		}
		data[i].Approvers = s.removeEnterpriseOwnersInApprovers(data[i].Approvers)
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

func (s *itemService) removeEnterpriseOwnersInApprovers(approvers []string) []string {
	if len(approvers) == 1 {
		return approvers
	}

	ownersArray := s.configManager.GetEnterpriseOwners()
	if len(ownersArray) == 0 {
		return approvers
	}
	ownersMap := make(map[string]bool)
	for _, owner := range ownersArray {
		ownersMap[owner] = true
	}

	var newApprovers []string
	for _, approver := range approvers {
		if !ownersMap[approver] {
			newApprovers = append(newApprovers, approver)
		}
	}

	return newApprovers
}
