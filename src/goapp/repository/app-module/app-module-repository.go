package appModule

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type applicationModuleRepository struct {
	*db.Database
}

func NewApplicationModuleRepository(db *db.Database) ApplicationModuleRepository {
	return &applicationModuleRepository{
		Database: db,
	}
}

func (r *applicationModuleRepository) GetApplicationModuleByIdAndApplicationId(applicationId string, applicationModuleId string) (*model.ApplicationModule, error) {
	rowApplicationModule, err := r.Query("PR_ApplicationModules_Select_ById_ApplicationId",
		sql.Named("ApplicationId", applicationId),
		sql.Named("ApplicationModuleId", applicationModuleId))

	if err != nil {
		return nil, err
	}
	defer rowApplicationModule.Close()

	applicationModule, err := r.RowsToMap(rowApplicationModule)
	if err != nil {
		return nil, err
	}

	if len(applicationModule) == 0 {
		return nil, nil
	}

	result := model.ApplicationModule{
		ApplicationName:       applicationModule[0]["ApplicationName"].(string),
		ApplicationModuleName: applicationModule[0]["ApplicationModuleName"].(string),
		Callbackurl:           applicationModule[0]["CallbackUrl"].(string),
		RequireRemarks:        applicationModule[0]["RequireRemarks"].(bool),
		ApproveText:           applicationModule[0]["ApproveText"].(string),
		RejectText:            applicationModule[0]["RejectText"].(string),
		AllowReassign:         applicationModule[0]["AllowReassign"].(bool),
		RequireAuthentication: applicationModule[0]["RequireAuthentication"].(bool),
	}

	if applicationModule[0]["ReassignCallbackUrl"] != nil {
		result.ReassignCallbackUrl = applicationModule[0]["ReassignCallbackUrl"].(string)
	}

	if applicationModule[0]["ExportUrl"] != nil {
		result.ExportUrl = applicationModule[0]["ExportUrl"].(string)
	}

	return &result, nil
}

func (r *applicationModuleRepository) GetAll() ([]model.ApplicationModule, error) {
	rowApplicationModules, err := r.Query("PR_ApplicationModules_Select")
	if err != nil {
		return nil, err
	}
	defer rowApplicationModules.Close()

	applicationModules, err := r.RowsToMap(rowApplicationModules)
	if err != nil {
		return nil, err
	}

	var result []model.ApplicationModule
	for _, v := range applicationModules {
		m := model.ApplicationModule{
			ApplicationModuleId:   v["Id"].(string),
			ApplicationModuleName: v["Name"].(string),
		}

		if v["CallbackUrl"] != nil {
			m.Callbackurl = v["CallbackUrl"].(string)
		}

		if v["AllowReassign"] != nil {
			m.AllowReassign = v["AllowReassign"].(bool)
		}

		if v["RequireRemarks"] != nil {
			m.RequireRemarks = v["RequireRemarks"].(bool)
		}

		if v["RequireAuthentication"] != nil {
			m.RequireAuthentication = v["RequireAuthentication"].(bool)
		}

		if v["ReassignCallbackUrl"] != nil {
			m.ReassignCallbackUrl = v["ReassignCallbackUrl"].(string)
		}

		if v["ExportUrl"] != nil {
			m.ExportUrl = v["ExportUrl"].(string)
		}

		result = append(result, m)
	}

	return result, nil
}

func (r *applicationModuleRepository) IsAuthRequired(applicationModuleId string) (bool, error) {
	rows, err := r.Query("PR_ApplicationModules_IsAuthRequired",
		sql.Named("ApplicationModuleId", applicationModuleId))

	if err != nil {
		return true, err
	}
	defer rows.Close()

	result, err := r.RowsToMap(rows)
	if err != nil {
		return true, err
	}

	if len(result) == 0 {
		return result[0]["RequireAuthentication"].(bool), nil
	} else {
		return true, nil
	}
}
