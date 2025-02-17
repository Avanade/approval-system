package legalConsultation

import (
	"main/model"
)

type LegalConsultationService interface {
	GetAllLegalConsulations(filterOptions model.FilterOptions, status int) ([]model.Item, error)
	GetLegalConsultants(token string) ([]model.Approver, error)
	GetLegalConsultationByItemId(id string) ([]model.LegalConsultation, error)
	InsertLegalConsultation(lc *model.LegalConsultation) error
}
