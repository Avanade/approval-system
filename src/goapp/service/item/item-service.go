package item

import (
	"main/model"
	repositoryItem "main/repository/item"
)

type itemService struct {
	repositoryItem repositoryItem.ItemRepository
}

func NewItemService(repositoryItem repositoryItem.ItemRepository) ItemService {
	return &itemService{
		repositoryItem: repositoryItem,
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

	result = model.Response{
		Data:  data,
		Total: total,
	}

	return result, nil
}
