package appModule

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type applicationModuleRepository struct {
	db.Database
}

func NewApplicationModuleRepository(db db.Database) ApplicationModuleRepository {
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
