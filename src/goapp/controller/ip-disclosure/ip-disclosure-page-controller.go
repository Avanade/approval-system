package ipdisclosure

import (
	"main/service"
	"net/http"
)

type ipDisclosurePageController struct {
	*service.Service
}

func NewIpDisclosurePageController(s *service.Service) IpDisclosurePageController {
	return &ipDisclosurePageController{
		Service: s,
	}
}

func (c *ipDisclosurePageController) IpDisclosureRequest(w http.ResponseWriter, r *http.Request) {
	user, err := c.Service.Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, d := c.Service.Template.UseTemplate("ipdisclosurerequest", r.URL.Path, *user, nil)

	err = t.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
