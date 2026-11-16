package handles

import (
	"html/template"
	"shop/repo"
	"shop/secret"
)

type Handler struct {
	repo   *repo.Repo
	secret *secret.Secret
	tpl    *template.Template
}

func New(repo *repo.Repo, secret *secret.Secret) *Handler {
	return &Handler{
		repo:   repo,
		secret: secret,
		tpl:    template.Must(template.ParseGlob("handles/templates/*.html")),
		//загрузка шаблонов html, template - пакет Go для работы с html,
		// ParseGlob - загружает файлы по указанному пути
	}
}
