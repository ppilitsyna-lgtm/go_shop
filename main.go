package main

import (
	"log"
	"net/http"
	"shop/handles"
	"shop/middle"
	"shop/repo"
	"shop/secret"
)

func main() {
	r := repo.New()
	tokens := secret.New()
	h := handles.New(r, tokens)
	m := middle.New(tokens)

	mux := http.NewServeMux() //создает маршрутизатор

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Страницы, доступные всем
	mux.Handle("GET /{$}", m.Auth("*", m.Panic(m.Route(h.Home))))
	mux.Handle("GET /login", m.Auth("*", m.Panic(m.Route(h.LoginPage))))
	mux.Handle("POST /login", m.Auth("*", m.Panic(m.Route(h.Login))))
	mux.Handle("GET /register", m.Auth("*", m.Panic(m.Route(h.RegisterPage))))
	mux.Handle("POST /register", m.Auth("*", m.Panic(m.Route(h.Register))))
	mux.Handle("GET /logout", m.Auth("*", m.Panic(m.Route(h.Logout))))
	//открывается главная страница -> страница входа -> обработка формы входа и тд
	//Route запись запросов в консоль

	// Страницы, доступные пользователям, которые вошли в аккаунт
	mux.Handle("GET /cart", m.Auth("user", m.Panic(m.Route(h.Cart))))
	mux.Handle("POST /cart/add/{id}", m.Auth("user", m.Panic(m.Route(h.CartAdd))))
	mux.Handle("POST /cart/remove/{id}", m.Auth("user", m.Panic(m.Route(h.CartRemove))))
	mux.Handle("POST /cart/checkout", m.Auth("user", m.Panic(m.Route(h.CartCheckout))))
	//просмотр корзины, добавление товара в корзину, удаление из корзины, оформление заказа
	// m.Auth для проверки авторизации и роли пользователя

	// Управление товарами, только для админа
	mux.Handle("GET /admin/products", m.Auth("admin", m.Panic(m.Route(h.ProductList))))
	mux.Handle("GET /admin/products/{id}", m.Auth("admin", m.Panic(m.Route(h.ProductGet))))
	mux.Handle("POST /admin/products/create", m.Auth("admin", m.Panic(m.Route(h.ProductCreate))))
	mux.Handle("POST /admin/products/update/{id}", m.Auth("admin", m.Panic(m.Route(h.ProductUpdate))))
	mux.Handle("POST /admin/products/delete/{id}", m.Auth("admin", m.Panic(m.Route(h.ProductDelete))))
	// Проверка авторизации и поли админа
	// Просмтор списка товаров, получение товара для редактирования, добавление, изменение, удаление товара

	// Управление категориями, только для админа
	mux.Handle("GET /admin/categories", m.Auth("admin", m.Panic(m.Route(h.CategoryList))))
	mux.Handle("GET /admin/categories/{id}", m.Auth("admin", m.Panic(m.Route(h.CategoryGet))))
	mux.Handle("POST /admin/categories/create", m.Auth("admin", m.Panic(m.Route(h.CategoryCreate))))
	mux.Handle("POST /admin/categories/update/{id}", m.Auth("admin", m.Panic(m.Route(h.CategoryUpdate))))
	mux.Handle("POST /admin/categories/delete/{id}", m.Auth("admin", m.Panic(m.Route(h.CategoryDelete))))
	// Просмотр списка категорий, получение категории для редактирования, добавление, изменение, удаление категории

	srv := &http.Server{Addr: ":8080", Handler: mux}
	log.Println("Сайт запущен: http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
