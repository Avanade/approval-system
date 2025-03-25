package fallback

import (
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
	user, err := c.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, d := c.Service.Template.UseTemplate("NotFound", r.URL.Path, *user, nil)
	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
