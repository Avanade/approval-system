package itemActivity

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type itemActivityRepository struct {
	*db.Database
}

func NewItemActivityRepository(db *db.Database) ItemActivityRepository {
	return &itemActivityRepository{
		Database: db,
	}
}

func (r *itemActivityRepository) GetItemActivityByItemId(id string) ([]model.ItemActivity, error) {
	var result []model.ItemActivity
	row, err := r.Query("PR_ItemActivity_Select_ByItemId", sql.Named("ItemId", id))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	itemActivity, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	for _, item := range itemActivity {
		result = append(result, model.ItemActivity{
			Id:        item["Id"].(int64),
			CreatedBy: item["CreatedBy"].(string),
			Created:   item["Created"].(time.Time).String(),
			Content:   item["Content"].(string),
		})
	}

	return result, nil
}

func (r *itemActivityRepository) AddItemActivity(itemActivity *model.ItemActivity) (*model.ItemActivity, error) {
	row, err := r.QueryRow("PR_ItemActivity_Insert",
		sql.Named("CreatedBy", itemActivity.CreatedBy),
		sql.Named("Content", itemActivity.Content),
		sql.Named("ItemId", itemActivity.ItemId),
	)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&itemActivity.Id)
	if err != nil {
		return nil, err
	}

	return itemActivity, nil
}
