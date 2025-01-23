package template

import (
	"fmt"
	"main/config"
	"main/model"
	"strings"
	"text/template"
)

type templateService struct {
	LinkFooters      string
	OrganizationName string
}

func NewTemplateService(config config.ConfigManager) TemplateService {
	return &templateService{
		LinkFooters:      config.GetLinkFooters(),
		OrganizationName: config.GetOrganizationName(),
	}
}

func (t *templateService) UseTemplate(page, path string, user model.AzureUser, pageData interface{}) (*template.Template, *model.MasterPageData) {
	// Data on master page
	var menu []model.Menu
	menu = append(menu, model.Menu{Name: "My Requests", Url: "/", IconPath: "/public/icons/projects.svg"})
	menu = append(menu, model.Menu{Name: "My Approvals", Url: "/myapprovals", IconPath: "/public/icons/approvals.svg"})
	menu = append(menu, model.Menu{Name: "Request IP Disclosure", Url: "/ipdisclosurerequest", IconPath: "/public/icons/ipdisclosure.svg"})
	masterPageData := model.Headers{Menu: menu, Page: getUrlPath(path)}

	//Footers
	var footers []model.Footer
	footerString := t.LinkFooters
	res := strings.Split(footerString, ";")
	for _, footer := range res {
		f := strings.Split(footer, ">")
		footers = append(footers, model.Footer{Text: f[0], Url: f[1]})
	}

	data := model.MasterPageData{
		Header:           masterPageData,
		Profile:          user,
		Content:          pageData,
		Footers:          footers,
		OrganizationName: t.OrganizationName,
	}

	tmpl := template.Must(
		template.ParseFiles("templates/master.html", "templates/buttons.html",
			fmt.Sprintf("templates/%v.html", page)))

	return tmpl, &data
}

func getUrlPath(path string) string {
	p := strings.Split(path, "/")
	if p[1] == "" {
		return "/"
	} else {
		return fmt.Sprintf("/%s", p[1])
	}
}
