package ipdrequest

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type ipdRequestRepository struct {
	*db.Database
}

func NewIpdRequestRepository(db *db.Database) IpdRequestRepository {
	return &ipdRequestRepository{
		Database: db,
	}
}

func (r *ipdRequestRepository) GetIpdRequestByApprovalRequestId(approvalRequestId string) (*model.IPDRequest, error) {
	row, err := r.Query("PR_IPDisclosureRequest_Select_ByApprovalRequestId",
		sql.Named("ApprovalRequestId", approvalRequestId),
	)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	ipdRequest := model.IPDRequest{
		RequestId:         result[0]["Id"].(int64),
		RequestorName:     result[0]["RequestorName"].(string),
		RequestorEmail:    result[0]["RequestorEmail"].(string),
		IPTitle:           result[0]["IPTitle"].(string),
		IPType:            result[0]["IPType"].(string),
		IPDescription:     result[0]["IPDescription"].(string),
		Reason:            result[0]["Reason"].(string),
		ApprovalRequestId: approvalRequestId,
		Created:           result[0]["Created"].(time.Time).String(),
	}

	if result[0]["ResponseDate"] != nil {
		ipdRequest.ResponseDate = result[0]["ResponseDate"].(time.Time).String()
	}

	if result[0]["IsApproved"] != nil {
		ipdRequest.IsApproved = result[0]["IsApproved"].(bool)
	}

	if result[0]["ApproverRemarks"] != nil {
		ipdRequest.ApproverRemarks = result[0]["ApproverRemarks"].(string)
	}

	if result[0]["RespondedBy"] != nil {
		ipdRequest.RespondedBy = result[0]["RespondedBy"].(string)
	}

	return &ipdRequest, nil
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

func (r *ipdRequestRepository) UpdateResponse(data *model.ResponseCallback) error {
	_, err := r.Query("PR_IPDisclosureRequest_Update_Response",
		sql.Named("ApprovalRequestId", data.ItemId),
		sql.Named("IsApproved", data.IsApproved),
		sql.Named("ApproverRemarks", data.Remarks),
		sql.Named("RespondedBy", data.RespondedBy),
	)

	return err
}
