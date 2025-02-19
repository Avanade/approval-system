package permission

type PermissionService interface {
	GetUserWithPermission(permission string) ([]string, error)
}
