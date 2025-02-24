package involvement

import (
	db "main/infrastructure/database"
	"main/model"
)

type involvementRepository struct {
	*db.Database
}

func NewInvolvementRepository(db *db.Database) InvolvementRepository {
	return &involvementRepository{
		Database: db,
	}
}

func (r *involvementRepository) GetInvolvementList() ([]model.Involvement, error) {
	var result []model.Involvement
	rows, err := r.Query("PR_Involvement_Select_All")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	involvements, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, involvement := range involvements {
		result = append(result, model.Involvement{
			Id:   involvement["Id"].(string),
			Name: involvement["Name"].(string),
		})
	}

	return result, nil
}
