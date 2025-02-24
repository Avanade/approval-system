package permission

import (
	"database/sql"
	db "main/infrastructure/database"
)

type permissionRepository struct {
	*db.Database
}

func NewPermissionRepository(db *db.Database) PermissionRepository {
	return &permissionRepository{
		Database: db,
	}
}

func (r *permissionRepository) GetUsersWithPermission(permission string) ([]string, error) {
	var users []string

	row, err := r.Query("PR_Permission_Select_ByType",
		sql.Named("Type", permission))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var user string
		err := row.Scan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
