package permission

type PermissionRepository interface {
	GetUsersWithPermission(permission string) ([]string, error)
}
