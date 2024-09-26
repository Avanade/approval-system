package service

import (
	"main/config"
	"main/repository"
	sApplicationModule "main/service/app-module"
	sApprovalRequestApprover "main/service/approval-request-approver"
	sEmail "main/service/email"
	sItem "main/service/item"
)

type Service struct {
	ApplicationModule       sApplicationModule.ApplicationModuleService
	Item                    sItem.ItemService
	Email                   sEmail.EmailService
	ApprovalRequestApprover sApprovalRequestApprover.ApprovalRequestApproverService
}

type ServiceOptionFunc func(*Service)

func NewService(opts ...ServiceOptionFunc) *Service {
	service := &Service{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func NewApplicationModuleService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ApplicationModule = sApplicationModule.NewApplicationModuleService(repo)
	}
}

func NewItemService(repo *repository.Repository, conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Item = sItem.NewItemService(repo, conf)
	}
}

func NewEmailService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Email = sEmail.NewSdkEmailService(conf)
	}
}

func NewApprovalRequestApproverService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ApprovalRequestApprover = sApprovalRequestApprover.NewApprovalRequestApproverService(repo)
	}
}
