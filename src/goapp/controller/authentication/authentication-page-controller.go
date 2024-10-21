package authentication

import (
	"html/template"
	"net/http"
)

type authenticationPageController struct{}

func NewAuthenticationController() AuthenticationPageController {
	return &authenticationPageController{}
}

func (a *authenticationPageController) LoginRedirectHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	redirect := "/"
	if len(q["redirect"]) > 0 {
		redirect = q["redirect"][0]
	}
	data := map[string]interface{}{
		"redirect": redirect,
	}

	c := http.Cookie{
		Name:   "auth-session",
		MaxAge: -1}
	http.SetCookie(w, &c)

	tmpl := template.Must(template.ParseFiles("templates/loginredirect.html"))
	tmpl.Execute(w, data)
}