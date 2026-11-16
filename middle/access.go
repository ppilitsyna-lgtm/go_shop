package middle

import (
	"net/http"
	"shop/secret"
	"strings"
)

func (m *Middle) Auth(acc string, next http.Handler) http.Handler { //необходимый ур-нь доступа, обработчик, который выполнится если доступ разрешен
	accessNum := access(acc) //ф-я преобразовывает название роли в число
	if accessNum == 1000 {
		panic("неизвестный уровень доступа: " + acc)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accessNum >= 100 { //проверка авторизации не выполняется(главная,вход,регистрация)
			next.ServeHTTP(w, r)
			return
		}
		usr, err := m.secret.Parse(secret.Token(r)) //если страница с авторизацией, получаем данные пользователя из токена
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if access(usr.Role) <= accessNum { //если значение роли пользователя <= ур-ню доступа страницы, доступ разрешен
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	})
}

func access(name string) int {
	accessNum := 1000
	switch strings.ToLower(name) { //все буквы в нижний регистр
	case "*", "all", "full":
		accessNum = 100 //страница доступна всем
	case "user", "users", "client":
		accessNum = 10 //страница доступна обычному пользователю
	case "admin", "admins", "root":
		accessNum = 1 //страница доступна только админу
	}
	return accessNum
}
