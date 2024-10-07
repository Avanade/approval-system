package model

type Application struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	ExportUrl           string `json:"exportUrl"`
	OrganizationTypeUrl string `json:"organizationTypeUrl"`
}
