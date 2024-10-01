package appmodule

import (
	"encoding/json"
	"main/service"
	"net/http"
)

type applicationModuleController struct {
	*service.Service
}

func NewApplicationModuleController(s *service.Service) ApplicationModuleController {
	return &applicationModuleController{
		Service: s,
	}
}

func (c *applicationModuleController) GetRequestTypes(w http.ResponseWriter, r *http.Request) {
	requestTypes, err := c.Service.ApplicationModule.GetRequestTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestTypes)
}
