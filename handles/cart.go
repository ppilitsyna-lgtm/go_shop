package handles

import (
	"net/http"
	"shop/domain"
	"strconv"
)

func (h *Handler) Cart(w http.ResponseWriter, r *http.Request) {
	usr := h.user(r)
	rows, err := h.repo.Db.QueryContext(r.Context(), `select p.id, p.name, p.description, p.price, p.image, c.qty
		from cart c join products p on p.id = c.product_id
		where c.user_id = $1 order by p.id`, usr.Id)
	if err != nil {
		h.fail(w, err)
		return
	}
	defer rows.Close()

	var items []domain.CartItem
	total := 0
	for rows.Next() {
		var it domain.CartItem
		p := &it.Product
		if err = rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Image, &it.Quantity); err != nil {
			h.fail(w, err)
			return
		}
		it.Sum = p.Price * it.Quantity
		total += it.Sum
		items = append(items, it)
	}
	if err = rows.Err(); err != nil {
		h.fail(w, err)
		return
	}
	h.render(w, r, "cart.html", map[string]any{
		"Items": items,
		"Total": total,
		"Done":  r.URL.Query().Get("done") == "1",
	})
}

func (h *Handler) CartAdd(w http.ResponseWriter, r *http.Request) {
	usr := h.user(r)
	id, _ := strconv.Atoi(r.PathValue("id"))
	_, err := h.repo.Db.ExecContext(r.Context(), `insert into cart (user_id, product_id, qty) values ($1, $2, 1)
		on conflict (user_id, product_id) do update set qty = cart.qty + 1`, usr.Id, id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) CartRemove(w http.ResponseWriter, r *http.Request) {
	usr := h.user(r)
	id, _ := strconv.Atoi(r.PathValue("id"))
	_, err := h.repo.Db.ExecContext(r.Context(), "delete from cart where user_id = $1 and product_id = $2", usr.Id, id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) CartCheckout(w http.ResponseWriter, r *http.Request) {
	usr := h.user(r)
	_, err := h.repo.Db.ExecContext(r.Context(), "delete from cart where user_id = $1", usr.Id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/cart?done=1", http.StatusSeeOther)
}
