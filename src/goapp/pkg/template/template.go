package template

import (
	"fmt"
	"html/template"
	"main/model"
	session "main/pkg/session"
	"net/http"
	"os"
	"strings"
)

// This parses the master page layout and the required page template.
func UseTemplate(w *http.ResponseWriter, r *http.Request, page string, pageData interface{}) error {

	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return err
	}
	profile := sessionaz.Values["profile"]

	user := model.AzureUser{}
	if profile != nil {
		user.Name = profile.(map[string]interface{})["name"].(string)
		user.Email = profile.(map[string]interface{})["preferred_username"].(string)
	}

	// Data on master page
	var menu []model.Menu
	menu = append(menu, model.Menu{Name: "My Requests", Url: "/", IconPath: "/public/icons/projects.svg"})
	menu = append(menu, model.Menu{Name: "My Approvals", Url: "/myapprovals", IconPath: "/public/icons/approvals.svg"})
	masterPageData := model.Headers{Menu: menu, Page: getUrlPath(r.URL.Path)}

	//Footers
	var footers []model.Footer
	footerString := os.Getenv("LINK_FOOTERS")
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
		OrganizationName: os.Getenv("ORGANIZATION_NAME"),
	}

	tmpl := template.Must(
		template.ParseFiles("templates/master.html", "templates/buttons.html",
			fmt.Sprintf("templates/%v.html", page)))

	return tmpl.Execute(*w, data)
}

func getUrlPath(path string) string {
	p := strings.Split(path, "/")
	if p[1] == "" {
		return "/"
	} else {
		return fmt.Sprintf("/%s", p[1])
	}
}
