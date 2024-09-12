package item

import (
	"main/model"
	repositoryItem "main/repository/item"
	"net/url"
	"strconv"
)

type itemService struct {
	repositoryItem repositoryItem.ItemRepository
}

func NewItemService(repositoryItem repositoryItem.ItemRepository) ItemService {
	return &itemService{
		repositoryItem: repositoryItem,
	}
}

func (s *itemService) GetAll(routeVars map[string]string, params url.Values, user string) (model.Response, error) {
	itemType, err := strconv.ParseInt(routeVars["type"], 10, 8)
	if err != nil {
		return model.Response{}, err
	}

	itemStatus, err := strconv.ParseInt(routeVars["status"], 10, 8)
	if err != nil {
		return model.Response{}, err
	}

	var result model.Response

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		requestType := ""
		if params["requestType"] != nil {
			if params["requestType"][0] != "" {
				requestType = params["requestType"][0]
			}
		}

		organization := ""
		if params["organization"] != nil {
			if params["organization"][0] != "" {
				organization = params["organization"][0]
			}
		}

		total, err := s.repositoryItem.GetTotalItemsBy(model.ItemType(itemType), model.ItemStatus(itemStatus), requestType, organization, user, search)
		if err != nil {
			return model.Response{}, err
		}

		data, err := s.repositoryItem.GetItemsBy(model.ItemType(itemType), model.ItemStatus(itemStatus), requestType, organization, user, search, offset, filter)
		if err != nil {
			return model.Response{}, err
		}

		result = model.Response{
			Data:  data,
			Total: total,
		}
	}
	return result, nil
}
