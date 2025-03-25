package permission

import (
	"main/repository"
)

type permissionService struct {
	Repository *repository.Repository
}

func NewPermissionService(repo *repository.Repository) PermissionService {
	return &permissionService{
		Repository: repo,
	}
}

func (s *permissionService) GetUserWithPermission(permission string) ([]string, error) {
	return s.Repository.Permission.GetUsersWithPermission(permission)
}
