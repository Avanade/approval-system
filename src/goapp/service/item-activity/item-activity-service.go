package itemActivity

import (
	"main/model"
	"main/repository"
)

type itemActivityService struct {
	Repository *repository.Repository
}

func NewItemActivityService(repo *repository.Repository) ItemActivityService {
	return &itemActivityService{
		Repository: repo,
	}
}

func (s *itemActivityService) GetItemActivityByItemId(id string) ([]model.ItemActivity, error) {
	return s.Repository.ItemActivity.GetItemActivityByItemId(id)
}

func (s *itemActivityService) AddItemActivity(itemActivity *model.ItemActivity) (*model.ItemActivity, error) {
	return s.Repository.ItemActivity.AddItemActivity(itemActivity)
}
