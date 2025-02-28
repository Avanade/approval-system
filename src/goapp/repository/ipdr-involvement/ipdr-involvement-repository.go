package ipdrInvolvement

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"strconv"
)

type ipdrInvolvementRepository struct {
	*db.Database
}

func NewIpdrInvolvementRepository(db *db.Database) IpdrInvolvementRepository {
	return &ipdrInvolvementRepository{
		Database: db,
	}
}

func (r *ipdrInvolvementRepository) GetIpdrInvolvementByRequestId(requestId int64) (involvementIds []string, involvements []string, err error) {
	row, err := r.Query("PR_IPDRInvolvement_Select_ByIPDRId",
		sql.Named("RequestId", requestId),
	)

	if err != nil {
		return nil, nil, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return nil, nil, err
	}

	for _, row := range result {
		involvementIds = append(involvementIds, strconv.FormatInt(row["InvolvementId"].(int64), 10))
		involvements = append(involvements, row["InvolvementName"].(string))
	}

	return involvementIds, involvements, nil
}

func (r *ipdrInvolvementRepository) InsertIpdrInvolvement(ipdrInvolvement model.IpdrInvolvement) error {
	involvementId, err := strconv.Atoi(ipdrInvolvement.InvolvementId)
	if err != nil {
		return err
	}
	row, err := r.Query("PR_IPDRInvolvement_Insert",
		sql.Named("RequestId", ipdrInvolvement.RequestId),
		sql.Named("InvolvementId", involvementId),
	)

	if err != nil {
		return err
	}
	defer row.Close()

	return nil
}
