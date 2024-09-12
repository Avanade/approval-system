package item

import (
	"main/model"
	"net/url"
)

type ItemService interface {
	GetAll(routeVars map[string]string, params url.Values, user string) (model.Response, error)
}
