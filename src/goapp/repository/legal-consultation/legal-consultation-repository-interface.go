package legalConsultation

import (
	"main/model"
)

type LegalConsultationRepository interface {
	GetLegalConsultationByEmail(email string, filterOptions model.FilterOptions, status int) ([]model.Item, error)
	GetLegalConsultationByItemId(itemId string) ([]model.LegalConsultation, error)
	GetTotalLegalConsultationByEmail(email string, status int) (int, error)
	InsertLegalConsultation(lc *model.LegalConsultation) error
}
