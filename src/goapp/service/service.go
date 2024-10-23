package service

import (
	"main/config"
	"main/infrastructure/session"
	"main/repository"
	sApplicationModule "main/service/app-module"
	sApplication "main/service/application"
	sApprovalRequestApprover "main/service/approval-request-approver"
	sAuthenticator "main/service/authenticator"
	sEmail "main/service/email"
	sItem "main/service/item"
	sMsGraph "main/service/msgraph"
	sTemplate "main/service/template"
)

type Service struct {
	Application             sApplication.ApplicationService
	ApplicationModule       sApplicationModule.ApplicationModuleService
	Item                    sItem.ItemService
	Email                   sEmail.EmailService
	ApprovalRequestApprover sApprovalRequestApprover.ApprovalRequestApproverService
	MsGraph                 sMsGraph.MsGraphService
	Template                sTemplate.TemplateService
	Authenticator           sAuthenticator.AuthenticatorService
}

type ServiceOptionFunc func(*Service)

func NewService(opts ...ServiceOptionFunc) *Service {
	service := &Service{}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func NewApplicationService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.Application = sApplication.NewApplicationService(repo)
	}
}

func NewApplicationModuleService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ApplicationModule = sApplicationModule.NewApplicationModuleService(repo)
	}
}

func NewAuthenticatorService(conf config.ConfigManager, session *session.Session) ServiceOptionFunc {
	return func(s *Service) {
		s.Authenticator = sAuthenticator.NewAuthenticatorService(conf, session)
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

func NewMsGraphService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.MsGraph = sMsGraph.NewMsGraphService(conf)
	}
}

func NewTemplateService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Template = sTemplate.NewTemplateService(conf)
	}
}
