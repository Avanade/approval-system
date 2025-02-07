package item

import (
	"fmt"
	"main/config"
	"main/model"
	"main/repository"
	"strconv"
	"sync"
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

func (s *itemService) GetFailedCallbacks() ([]string, error) {
	failedCallbacks, err := s.Repository.Item.GetFailedCallbacks()
	if err != nil {
		return []string{}, err
	}
	return failedCallbacks, nil
}

func (s *itemService) GetItemById(id string) (*model.Item, error) {
	item, err := s.Repository.Item.GetItemById(id)
	if err != nil {
		return nil, err
	}
	return item, nil
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

	var wg sync.WaitGroup
	maxGoroutines := 2
	guard := make(chan struct{}, maxGoroutines)

	for i := range data {
		guard <- struct{}{}
		wg.Add(1)
		go func(r *model.Item) {
			approvers, err := s.Repository.ApprovalRequestApprover.GetApproversByItemId(r.Id)
			if err != nil {
				fmt.Println("Error getting approvers of item id: ", r.Id)
				return
			}
			r.Approvers = approvers

			if len(r.Approvers) > 0 {
				r.Approvers = s.removeEnterpriseOwnersInApprovers(r.Approvers)
			}

			<-guard
			wg.Done()
		}(&data[i])
	}
	wg.Wait()

	result = model.Response{
		Data:  data,
		Total: total,
	}

	return result, nil
}

func (s *itemService) GetByApprover(approver, requestType, organization string, filterOptions model.FilterOptions) (items []model.Item, total int, err error) {
	items, total, err = s.Repository.Item.GetItemsByApprover(approver, requestType, organization, filterOptions)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *itemService) GetInvolvedUsers(itemId string) ([]string, error) {
	var users []string
	approvers, err := s.Repository.ApprovalRequestApprover.GetApproversByItemId(itemId)
	if err != nil {
		return nil, err
	}

	for _, approver := range approvers {
		users = append(users, approver)
	}

	item, err := s.Repository.Item.GetItemById(itemId)
	if err != nil {
		return nil, err
	}

	users = append(users, item.RequestedBy)

	return users, nil
}

func (s *itemService) InsertItem(item model.ItemInsertRequest) (string, error) {
	id, err := s.Repository.Item.InsertItem(item.ApplicationModuleId, item.Subject, item.Body, item.RequesterEmail)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *itemService) ItemIsAuthorized(appId, appModuleId, itemId, approverEmail string) (*model.ItemIsAuthorized, error) {
	itemIsAuthorized, err := s.Repository.Item.ItemIsAuthorized(appId, appModuleId, itemId, approverEmail)
	if err != nil {
		return nil, err
	}
	return itemIsAuthorized, nil
}

func (s *itemService) UpdateItemApproverEmail(itemId, approverEmail, username string) error {
	err := s.Repository.Item.UpdateItemApproverEmail(itemId, approverEmail, username)
	if err != nil {
		return err
	}
	return nil
}

func (s *itemService) UpdateItemCallback(itemId string, isCallbackFailed bool) error {
	err := s.Repository.Item.UpdateItemCallback(itemId, isCallbackFailed)
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

func (s *itemService) UpdateItemResponse(req model.ProcessResponseRequest) error {
	isApproved, err := strconv.ParseBool(req.IsApproved)
	if err != nil {
		return err
	}

	err = s.Repository.Item.UpdateItemResponse(req.ItemId, req.Remarks, req.ApproverEmail, isApproved)
	if err != nil {
		return err
	}
	return nil
}

func (s *itemService) ValidateItem(req model.ProcessResponseRequest) (bool, error) {
	isValid, err := s.Repository.Item.ValidateItem(req.ApplicationId, req.ApplicationModuleId, req.ItemId, req.ApproverEmail)
	if err != nil {
		return false, err
	}
	return isValid, nil
}
