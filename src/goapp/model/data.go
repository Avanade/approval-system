package model

type TypPageData struct {
	Header           interface{}
	Profile          interface{}
	Content          interface{}
	Footers          []Footer
	OrganizationName string
}

type Footer struct {
	Text string
	Url  string
}

type TypHeaders struct {
	Menu          []TypMenu
	ExternalLinks []TypMenu
	Page          string
}

type TypMenu struct {
	Name     string
	Url      string
	IconPath string
	UrlPath  string
}

type TypAzureUser struct {
	Name  string `json:"name"`
	Email string `json:"preferred_username"`
}
