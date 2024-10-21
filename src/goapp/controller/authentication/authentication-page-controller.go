package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	session "main/pkg/session"
	"main/service"
	"net/http"
)

type authenticationPageController struct {
	*service.Service
}

func NewAuthenticationController(s *service.Service) AuthenticationPageController {
	return &authenticationPageController{
		Service: s,
	}
}

func (a *authenticationPageController) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check session
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	//Retrieve token and save data on session store
	u, err := a.Service.Authenticator.ProcessToken(r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = u.IdToken
	session.Values["access_token"] = u.AccessToken
	session.Values["profile"] = u.Profile
	session.Values["refresh_token"] = u.RefreshToken
	session.Values["expiry"] = u.Expiry
	session.Options.MaxAge = 43200
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to index
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *authenticationPageController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := session.Store.Get(r, "auth-session")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, a.Authenticator.GetAuthCodeURL(state), http.StatusTemporaryRedirect)
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
