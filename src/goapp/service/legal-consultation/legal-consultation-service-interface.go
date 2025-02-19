package legalConsultation

import (
	"main/model"
)

type LegalConsultationService interface {
	GetLegalConsultants(token string) ([]model.Approver, error)
	GetLegalConsultationByItemId(id string) ([]model.LegalConsultation, error)
	InsertLegalConsultation(lc *model.LegalConsultation) error
}
