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
	sInvolvement "main/service/involvement"
	sIPDisclosureRequest "main/service/ip-disclosure-request"
	sItem "main/service/item"
	sItemActivity "main/service/item-activity"
	sLegalConsultation "main/service/legal-consultation"
	sMsGraph "main/service/msgraph"
	sPermission "main/service/permission"
	sTemplate "main/service/template"
)

type Service struct {
	Application             sApplication.ApplicationService
	ApplicationModule       sApplicationModule.ApplicationModuleService
	ApprovalRequestApprover sApprovalRequestApprover.ApprovalRequestApproverService
	Authenticator           sAuthenticator.AuthenticatorService
	Email                   sEmail.EmailService
	Involvement             sInvolvement.InvolvementService
	IPDisclosureRequest     sIPDisclosureRequest.IpDisclosureRequestService
	Item                    sItem.ItemService
	ItemActivity            sItemActivity.ItemActivityService
	LegalConsultation       sLegalConsultation.LegalConsultationService
	MsGraph                 sMsGraph.MsGraphService
	Permission              sPermission.PermissionService
	Template                sTemplate.TemplateService
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

func NewAuthenticatorService(conf config.ConfigManager, session *session.ConnectSession) ServiceOptionFunc {
	return func(s *Service) {
		s.Authenticator = sAuthenticator.NewAuthenticatorService(conf, *session)
	}
}

func NewInvolvementService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.Involvement = sInvolvement.NewInvolvementService(repo)
	}
}

func NewIPDisclosureRequestService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.IPDisclosureRequest = sIPDisclosureRequest.NewIpDisclosureRequestService(repo)
	}
}

func NewItemService(repo *repository.Repository, conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Item = sItem.NewItemService(repo, conf)
	}
}

func NewItemActivityService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ItemActivity = sItemActivity.NewItemActivityService(repo)
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

func NewLegalConsultationService(repo *repository.Repository, conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.LegalConsultation = sLegalConsultation.NewLegalConsultationService(repo, conf)
	}
}

func NewMsGraphService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.MsGraph = sMsGraph.NewMsGraphService(conf)
	}
}

func NewPermissionService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.Permission = sPermission.NewPermissionService(repo)
	}
}

func NewTemplateService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Template = sTemplate.NewTemplateService(conf)
	}
}
