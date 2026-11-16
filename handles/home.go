package handles

import (
	"net/http"
	"strconv" //преобразование строк в числа и обратно
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	cat, _ := strconv.Atoi(r.URL.Query().Get("cat"))
	products, err := h.products(r.Context(), cat)
	if err != nil {
		h.fail(w, err)
		return
	}
	categories, err := h.categories(r.Context())
	if err != nil {
		h.fail(w, err)
		return
	}
	h.render(w, r, "home.html", map[string]any{
		"Products":   products,
		"Categories": categories,
		//функция отображения страницы, название html, набор данных
	})
}
