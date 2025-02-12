package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
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

	// Pull list of legal approvers by using the endpoint /api/repository-approvers/legal
	token, err := a.Service.Authenticator.GenerateToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	legalApprovers, err := a.Service.LegalConsultation.GetLegalConsultants(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is a legal approver
	isLegalApprover := false
	for _, v := range legalApprovers {
		if v.ApproverEmail == u.Profile["preferred_username"].(string) {
			isLegalApprover = true
			break
		}
	}

	data := map[string]interface{}{
		"id_token":        u.IdToken,
		"access":          u.AccessToken,
		"profile":         u.Profile,
		"refresh_token":   u.RefreshToken,
		"expiry":          u.Expiry,
		"isLegalApprover": isLegalApprover,
	}

	err = a.Authenticator.SaveOnSession(&w, r, "auth-session", data)

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

	data := map[string]interface{}{
		"state": state,
	}

	err = a.Authenticator.SaveOnSession(&w, r, "auth-session", data)

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
