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
