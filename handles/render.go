package handles

import (
	"log"
	"net/http"
	"shop/domain"
	"shop/secret"
)

func (h *Handler) user(r *http.Request) domain.User {
	usr, _ := h.secret.Parse(secret.Token(r))
	return usr
}

//метод user, получаем запрос пользователя, ф-я возвращает пользователя из пакета domain
//получает токен авторизации из запроса, проверяет и расшифровывает токен, чтобы получить данные пользователя

func (h *Handler) render(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	if usr, err := h.secret.Parse(secret.Token(r)); err == nil {
		data["User"] = usr
	}
	//получает токен авторизации из запроса, проверяет и расшифровывает токен, чтобы получить данные пользователя
	//если ошибки нет, добавляем пользователя в данные для html

	if err := h.tpl.ExecuteTemplate(w, page, data); err != nil {
		log.Println(err)
		//выбираем нужный шаблон html и заполняем его данными, отправляем результат пользователю
		//название шаблона, инфо в шаблон, если страница не отобразилась то записываем ее в журнал
	}
}

//отправляем страницу пользователю, запрос пользователя, название html, данные которые будут переданы в html

func (h *Handler) fail(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
}

//отправляем ответ пользователю, полученная ошибка,
// записываем инфу об ошибке в журнал, отправляем соо об ошибке юзеру,
//текст пользователю, код 500 означающий внутреннюю ошибку сервера
