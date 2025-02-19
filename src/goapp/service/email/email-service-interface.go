package email

import "main/model"

type ContentType int

const (
	Html ContentType = iota
	Text
)

type EmailService interface {
	SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSetItem bool) error
	SendApprovalRequestEmail(req *model.ItemInsertRequest, appModule *model.ApplicationModule, id string) error
	SendActivityEmail(req *model.ItemActivity, recipients []string, domain, action string) error
	SendLegalConsultationRequestEmail(req *model.ConsultLegalRequest, user *model.AzureUser, domain string, recipients []string) error
}
