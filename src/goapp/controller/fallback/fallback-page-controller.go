package fallback

import (
	"main/model"
	"main/pkg/session"
	"main/service"
	"net/http"
)

type fallbackController struct {
	*service.Service
}

func NewFallbackController(s *service.Service) FallbackController {
	return &fallbackController{
		Service: s,
	}
}

func (c *fallbackController) NotFound(w http.ResponseWriter, r *http.Request) {
	session, _ := session.Store.Get(r, "auth-session")

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}
	user := model.AzureUser{
		Name:  profile["name"].(string),
		Email: profile["preferred_username"].(string),
	}

	t, d := c.Service.Template.UseTemplate("NotFound", r.URL.Path, user, nil)
	err := t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
