package involvement

import (
	"encoding/json"
	"fmt"
	"main/service"
	"net/http"
)

type involvementController struct {
	*service.Service
}

func NewInvolvementController(s *service.Service) InvolvementController {
	return &involvementController{
		Service: s,
	}
}

func (c *involvementController) GetInvolvementList(w http.ResponseWriter, r *http.Request) {
	involvements, err := c.Service.Involvement.GetInvolvementList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(involvements)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(b)
}
