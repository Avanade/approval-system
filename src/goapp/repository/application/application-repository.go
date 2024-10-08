package application

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type applicationRepository struct {
	*db.Database
}

func NewApplicationRepository(db *db.Database) ApplicationRepository {
	return &applicationRepository{
		Database: db,
	}
}

func (r *applicationRepository) GetApplicationById(id string) (*model.Application, error) {
	var application model.Application
	rowApplication, err := r.Query("PR_Applications_Select_ById", sql.Named("Id", id))
	if err != nil {
		return nil, err
	}

	applications, err := r.RowsToMap(rowApplication)
	if err != nil {
		return nil, err
	}

	if len(applications) == 0 {
		return nil, nil
	} else {
		application.Id = applications[0]["Id"].(string)
		application.Name = applications[0]["Name"].(string)
		if applications[0]["ExportUrl"] != nil {
			application.ExportUrl = applications[0]["ExportUrl"].(string)
		}
		if applications[0]["OrganizationTypeUrl"] != nil {
			application.OrganizationTypeUrl = applications[0]["OrganizationTypeUrl"].(string)
		}
	}
	return &application, nil
}
