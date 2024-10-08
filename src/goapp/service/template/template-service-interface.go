package template

import (
	"main/model"
	"text/template"
)

type TemplateService interface {
	UseTemplate(page, path string, user model.AzureUser, pageData interface{}) (*template.Template, *model.MasterPageData)
}
