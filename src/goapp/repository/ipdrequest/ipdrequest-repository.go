package ipdrequest

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type ipdRequestRepository struct {
	*db.Database
}

func NewIpdRequestRepository(db *db.Database) IpdRequestRepository {
	return &ipdRequestRepository{
		Database: db,
	}
}

func (r *ipdRequestRepository) InsertIpdRequest(ipdRequest *model.IPDRequest) (int64, error) {
	row, err := r.Query("PR_IPDisclosureRequest_Insert",
		sql.Named("RequestorName", ipdRequest.RequestorName),
		sql.Named("RequestorEmail", ipdRequest.RequestorEmail),
		sql.Named("IPTitle", ipdRequest.IPTitle),
		sql.Named("IPType", ipdRequest.IPType),
		sql.Named("IPDescription", ipdRequest.IPDescription),
		sql.Named("Reason", ipdRequest.Reason),
	)

	if err != nil {
		return 0, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return 0, err
	}

	return result[0]["Id"].(int64), nil

}

func (r *ipdRequestRepository) UpdateApprovalRequestId(approvalRequestId string, IPDRequestId int64) error {
	_, err := r.Query("PR_IPDisclosureRequest_Update_ApprovalRequestId",
		sql.Named("ApprovalRequestId", approvalRequestId),
		sql.Named("IPDRequestId", IPDRequestId),
	)

	return err
}
