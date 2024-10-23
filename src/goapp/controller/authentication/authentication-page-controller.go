package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"main/model"
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
	state, err := a.Authenticator.GetStringValue(r, "auth-session", "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	//Retrieve token and save data on session store
	u, err := a.Authenticator.ProcessToken(r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.Authenticator.SaveOnSession(&w, r, "auth-session",
		model.SessionStringValue{Key: "id_token", Value: u.IdToken},
		model.SessionStringValue{Key: "access", Value: u.AccessToken},
		model.SessionMapValue{Key: "profile", Value: u.Profile},
		model.SessionStringValue{Key: "refresh_token", Value: u.RefreshToken},
		model.SessionStringValue{Key: "expiry", Value: u.Expiry})

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

	err = a.Authenticator.SaveOnSession(&w, r, "auth-session",
		model.SessionStringValue{Key: "state", Value: state})

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

func (a *authenticationPageController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := a.Authenticator.ClearFromSession(&w, r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	url, err := a.Authenticator.GetLogoutURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
