package legalConsultation

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type legalConsultationRepository struct {
	*db.Database
}

func NewLegalConsultationRepository(db *db.Database) LegalConsultationRepository {
	return &legalConsultationRepository{
		Database: db,
	}
}

func (r *legalConsultationRepository) GetLegalConsultationByEmail(email string, filterOptions model.FilterOptions, status int) ([]model.Item, error) {
	var items []model.Item
	offset := filterOptions.Page
	if filterOptions.Page != 0 {
		offset = filterOptions.Page * filterOptions.Filter
	}

	row, err := r.Query("PR_LegalConsultation_Select_ByEmail",
		sql.Named("Email", email),
		sql.Named("Offset", offset),
		sql.Named("Filter", filterOptions.Filter),
		sql.Named("IsApproved", status))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	result, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		item := model.Item{
			Id:          v["ItemId"].(string),
			Application: v["Application"].(string),
			Module:      v["Module"].(string),
			Created:     v["Created"].(time.Time).String(),
			RequestedBy: v["CreatedBy"].(string),
		}

		if v["ApproverRemarks"] != nil {
			item.ApproverRemarks = v["ApproverRemarks"].(string)
		}

		if v["Body"] != nil {
			item.Body = v["Body"].(string)
		}

		if v["DateResponded"] != nil {
			item.DateResponded = v["DateResponded"].(time.Time).Format("2006-01-02T15:04:05.000Z")
		}

		if v["DateSent"] != nil {
			item.DateSent = v["DateSent"].(time.Time).String()
		}

		if v["IsApproved"] != nil {
			item.IsApproved = v["IsApproved"].(bool)
		}

		if v["Subject"] != nil {
			item.Subject = v["Subject"].(string)
		}

		if v["RespondedBy"] != nil {
			item.RespondedBy = v["RespondedBy"].(string)
		}

		if v["ApplicationModuleId"] != nil {
			item.ModuleId = v["ApplicationModuleId"].(string)
		}

		if v["ApplicationId"] != nil {
			item.ApplicationId = v["ApplicationId"].(string)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *legalConsultationRepository) GetLegalConsultationByItemId(itemId string) ([]model.LegalConsultation, error) {
	var result []model.LegalConsultation
	row, err := r.Query("PR_LegalConsultation_Select_ByItemId", sql.Named("ItemId", itemId))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	legalConsultation, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	for _, item := range legalConsultation {
		result = append(result, model.LegalConsultation{
			ItemId:    itemId,
			Email:     item["Email"].(string),
			Created:   item["Created"].(time.Time).String(),
			CreatedBy: item["CreatedBy"].(string),
		})
	}

	return result, nil
}

func (r *legalConsultationRepository) GetTotalLegalConsultationByEmail(email string, status int) (int, error) {
	row, err := r.QueryRow("PR_LegalConsultation_Select_TotalByEmail",
		sql.Named("Email", email),
		sql.Named("IsApproved", status),
	)

	if err != nil {
		return 0, err
	}

	var total int
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *legalConsultationRepository) InsertLegalConsultation(lc *model.LegalConsultation) error {
	_, err := r.QueryRow("PR_LegalConsultation_Insert",
		sql.Named("ItemId", lc.ItemId),
		sql.Named("Email", lc.Email),
		sql.Named("CreatedBy", lc.CreatedBy),
	)

	if err != nil {
		return err
	}

	return nil
}
