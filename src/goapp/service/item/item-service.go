package item

import (
	"main/config"
	"main/model"
	repositoryItem "main/repository/item"
)

type itemService struct {
	repositoryItem repositoryItem.ItemRepository
	configManager  config.ConfigManager
}

func NewItemService(repositoryItem repositoryItem.ItemRepository, configManager config.ConfigManager) ItemService {
	return &itemService{
		repositoryItem: repositoryItem,
		configManager:  configManager,
	}
}

func (s *itemService) GetAll(itemOptions model.ItemOptions) (model.Response, error) {
	var result model.Response

	total, err := s.repositoryItem.GetTotalItemsBy(itemOptions)
	if err != nil {
		return model.Response{}, err
	}

	data, err := s.repositoryItem.GetItemsBy(itemOptions)
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
