package legalConsultation

import (
	"encoding/json"
	"main/config"
	"main/model"
	"main/repository"
	"net/http"
	"time"
)

type legalConsultationService struct {
	Repository *repository.Repository
	Config     config.ConfigManager
}

func NewLegalConsultationService(repo *repository.Repository, conf config.ConfigManager) LegalConsultationService {
	return &legalConsultationService{
		Repository: repo,
		Config:     conf,
	}
}

func (s *legalConsultationService) GetLegalConsultants(token string) ([]model.Approver, error) {
	url := s.Config.GetCommunityPortalDomain() + "/api/repository-approvers/legal"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: time.Second * 90,
	}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	// Decode response
	var legalApprovers []model.Approver
	err = json.NewDecoder(res.Body).Decode(&legalApprovers)
	if err != nil {
		return nil, err
	}

	return legalApprovers, nil
}

func (s *legalConsultationService) GetLegalConsultationByItemId(id string) ([]model.LegalConsultation, error) {
	return s.Repository.LegalConsultation.GetLegalConsultationByItemId(id)
}

func (s *legalConsultationService) InsertLegalConsultation(lc *model.LegalConsultation) error {
	return s.Repository.LegalConsultation.InsertLegalConsultation(lc)
}
