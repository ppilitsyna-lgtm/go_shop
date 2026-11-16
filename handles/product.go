package handles

import (
	"context"
	"net/http"
	"shop/domain"
	"strconv"
)

const productSelect = `select p.id, p.name, p.description, p.price, p.image,
	coalesce(p.category_id, 0), coalesce(c.name, '')
	from products p left join categories c on c.id = p.category_id `

//если у товара нет категории, возвращается 0; если нет названия категории, возвращается пустая строка
//получаем данные из products, присоединяем таблицу категорий, связываем товары и категории по их ID

// контекст запроса, ID категории чтоб отфильтровать, ф-я возвращает список товаров или ошибку
func (h *Handler) products(ctx context.Context, cat int) ([]domain.Product, error) {
	rows, err := h.repo.Db.QueryContext(ctx,
		productSelect+"where $1 = 0 or p.category_id = $1 order by p.id desc", cat)
	if err != nil {
		return nil, err
	}
	//если категория 0 - показываем все, иначе только категорию товаров, сортировка сначала новые товары
	defer rows.Close()
	var list []domain.Product
	for rows.Next() {
		var p domain.Product
		err = rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Image, &p.CategoryId, &p.CategoryName)
		//для каждой строки создается p куда записываются товары, Scan переносит значения из столбцов БД в поля структуры продукт
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (h *Handler) ProductList(w http.ResponseWriter, r *http.Request) {
	h.productPage(w, r, domain.Product{})
}

//вызывает ф-ю  отображения страницы товара

func (h *Handler) ProductGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var p domain.Product
	err := h.repo.Db.QueryRowContext(r.Context(), productSelect+"where p.id = $1", id).
		Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Image, &p.CategoryId, &p.CategoryName)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.productPage(w, r, p)
}

func (h *Handler) ProductCreate(w http.ResponseWriter, r *http.Request) {
	p := productFromForm(r)
	_, err := h.repo.Db.ExecContext(r.Context(),
		"insert into products (name, description, price, image, category_id) values ($1, $2, $3, $4, nullif($5, 0))",
		p.Name, p.Description, p.Price, p.Image, p.CategoryId)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (h *Handler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	p := productFromForm(r)
	p.Id, _ = strconv.Atoi(r.PathValue("id"))
	_, err := h.repo.Db.ExecContext(r.Context(),
		"update products set name = $1, description = $2, price = $3, image = $4, category_id = nullif($5, 0) where id = $6",
		p.Name, p.Description, p.Price, p.Image, p.CategoryId, p.Id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (h *Handler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	_, err := h.repo.Db.ExecContext(r.Context(), "delete from products where id = $1", id)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (h *Handler) productPage(w http.ResponseWriter, r *http.Request, edit domain.Product) {
	products, err := h.products(r.Context(), 0)
	if err != nil {
		h.fail(w, err)
		return
	}
	categories, err := h.categories(r.Context())
	if err != nil {
		h.fail(w, err)
		return
	}
	h.render(w, r, "admin_products.html", map[string]any{
		"Products":   products,
		"Categories": categories,
		"Edit":       edit,
	})
}

func productFromForm(r *http.Request) domain.Product {
	price, _ := strconv.Atoi(r.FormValue("price"))
	cat, _ := strconv.Atoi(r.FormValue("category_id"))
	return domain.Product{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		Image:       r.FormValue("image"),
		CategoryId:  cat,
	}
}
