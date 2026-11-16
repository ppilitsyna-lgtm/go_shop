package handles

import (
	"net/http"
	"shop/domain"
	"shop/secret"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "login.html", nil)
}

//открывает страницу авторизации

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var usr domain.User
	err := h.repo.Db.QueryRowContext(r.Context(),
		"select id, login, password, role from users where login = $1", r.FormValue("login")).
		Scan(&usr.Id, &usr.Login, &usr.Password, &usr.Role)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(r.FormValue("password"))) != nil {
		h.render(w, r, "login.html", map[string]any{"Error": "Неверный логин или пароль"})
		return
	}
	h.startSession(w, r, &usr)
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "register.html", nil)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	usr := domain.User{Login: r.FormValue("login"), Role: "user"}
	pass := r.FormValue("password")
	if usr.Login == "" || len(pass) < 4 {
		h.render(w, r, "register.html", map[string]any{"Error": "Введите логин и пароль (от 4 символов)"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		h.fail(w, err)
		return
	}
	//преобразует пароль в массив байтов, получаем хеш

	err = h.repo.Db.QueryRowContext(r.Context(),
		"insert into users (login, password, role) values ($1, $2, $3) returning id",
		usr.Login, string(hash), usr.Role).Scan(&usr.Id)
	if err != nil {
		h.render(w, r, "register.html", map[string]any{"Error": "Такой логин уже занят"})
		return
	}
	h.startSession(w, r, &usr)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.secret.Delete(secret.Token(r))
	http.SetCookie(w, &http.Cookie{Name: "token", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) startSession(w http.ResponseWriter, r *http.Request, usr *domain.User) {
	token, err := h.secret.Issue(usr)
	if err != nil {
		h.fail(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", HttpOnly: true})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
