package model

type MasterPageData struct {
	Header           Headers
	Profile          AzureUser
	Content          interface{}
	Footers          []Footer
	OrganizationName string
}

type Footer struct {
	Text string
	Url  string
}

type Headers struct {
	Menu          []Menu
	ExternalLinks []Menu
	Page          string
}

type Menu struct {
	Name     string
	Url      string
	IconPath string
	UrlPath  string
}

type AzureUser struct {
	Name  string `json:"name"`
	Email string `json:"preferred_username"`
}
