package routes

import (
	"main/models"
	template "main/pkg/template"
	"net/http"
)

func GithubHandler(w http.ResponseWriter, r *http.Request, data *models.TypPageData) {
	template.UseTemplate(&w, data, "github")
}
