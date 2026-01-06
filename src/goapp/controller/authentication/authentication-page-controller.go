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

func (a *authenticationPageController) AuthenticationSuccessfulHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/authenticationsuccessful.html"))
	tmpl.Execute(w, nil)
}

func (a *authenticationPageController) AuthenticationInProgressHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/authenticationinprogress.html"))
	tmpl.Execute(w, nil)
}

func (a *authenticationPageController) AuthenticationFailedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/authenticationfailed.html"))
	tmpl.Execute(w, nil)
}

func (a *authenticationPageController) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	state, err := a.Authenticator.GetStringValue(r, "auth-session", "state")
	if err != nil {
		http.Redirect(w, r, "/login/azure", http.StatusSeeOther)
		return
	}

	if r.URL.Query().Get("state") != state {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	//Retrieve token and save data on session store
	u, err := a.Authenticator.ProcessToken(r.URL.Query().Get("code"))
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	// Pull list of legal approvers by using the endpoint /api/repository-approvers/legal
	token, err := a.Service.Authenticator.GenerateToken()
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	legalApprovers, err := a.Service.LegalConsultation.GetLegalConsultants(token)
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
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

	// Get List of users with "audit" permission
	auditors, err := a.Service.Permission.GetUserWithPermission("audit")
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	// Check if user is in the audit list
	isAuditor := false
	for _, v := range auditors {
		if v == u.Profile["preferred_username"].(string) {
			isAuditor = true
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
		"isAuditor":       isAuditor,
	}

	err = a.Authenticator.SaveOnSession(&w, r, "auth-session", data)

	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	// Redirect to index
	http.Redirect(w, r, "/authentication/azure/successful", http.StatusSeeOther)
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
	_ = a.Authenticator.ClearFromSession(&w, r, "auth-session")

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

	tmpl := template.Must(template.ParseFiles("templates/loginredirect.html"))
	tmpl.Execute(w, data)
}

func (a *authenticationPageController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	url, err := a.Authenticator.GetLogoutURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = a.Authenticator.ClearFromSession(&w, r, "auth-session")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
