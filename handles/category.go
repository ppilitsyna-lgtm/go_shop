package handles

import (
	"context"
	"net/http"
	"strconv"

	"shop/domain"
)

func (h *Handler) categories(ctx context.Context) ([]domain.Category, error) {
	rows, err := h.repo.Db.QueryContext(ctx, "select id, name from categories order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.Category
	for rows.Next() {
		var c domain.Category
		if err = rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (h *Handler) CategoryList(w http.ResponseWriter, r *http.Request) {
	h.categoryPage(w, r, domain.Category{})
}

func (h *Handler) CategoryGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var c domain.Category
	err := h.repo.Db.QueryRowContext(r.Context(), "select id, name from categories where id = $1", id).
		Scan(&c.Id, &c.Name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.categoryPage(w, r, c)
}

func (h *Handler) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	_, err := h.repo.Db.ExecContext(r.Context(), "insert into categories (name) values ($1)", r.FormValue("name"))
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *Handler) CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	_, err := h.repo.Db.ExecContext(r.Context(), "update categories set name = $1 where id = $2", r.FormValue("name"), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *Handler) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	_, err := h.repo.Db.ExecContext(r.Context(), "delete from categories where id = $1", id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *Handler) categoryPage(w http.ResponseWriter, r *http.Request, edit domain.Category) {
	categories, err := h.categories(r.Context())
	if err != nil {
		h.fail(w, err)
		return
	}
	h.render(w, r, "admin_categories.html", map[string]any{"Categories": categories, "Edit": edit})
}
