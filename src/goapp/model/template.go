package model

type MasterPageData struct {
	Header           Headers
	Profile          AzureUser
	Content          interface{} `json:"content"`
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
