package item

import (
	"encoding/json"
	"fmt"
	"main/pkg/session"
	service "main/service/item"
	"net/http"

	"github.com/gorilla/mux"
)

type itemController struct {
	itemService service.ItemService
}

func NewItemController(itemService service.ItemService) ItemController {
	return &itemController{
		itemService: itemService,
	}
}

// GetItems is a function to get all items
func (c *itemController) GetItems(w http.ResponseWriter, r *http.Request) {
	// Get all items
	vars := mux.Vars(r)

	params := r.URL.Query()

	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	user := fmt.Sprintf("%s", profile["preferred_username"])

	result, err := c.itemService.GetAll(vars, params, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
