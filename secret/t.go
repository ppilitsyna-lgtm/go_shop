package secret

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"shop/domain"
	"sync"
)

type Secret struct { //структура отвечающая за
	mu     sync.Mutex             //механизм блокировки данных
	tokens map[string]domain.User //карта, где хранятся токены и соответствующие им пользователи
}

func New() *Secret { //создается хранилище токенов
	return &Secret{tokens: map[string]domain.User{}}
}

func (s *Secret) Issue(usr *domain.User) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil { //массив заполняется случайными данными
		return "", err
	}
	token := hex.EncodeToString(b) //преобразует массив байтов в строку
	s.mu.Lock()
	s.tokens[token] = domain.User{Id: usr.Id, Login: usr.Login, Role: usr.Role} //сохранение пользователя
	s.mu.Unlock()
	return token, nil
}

func (s *Secret) Parse(token string) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	usr, ok := s.tokens[token] //поиск токена
	if !ok {
		return domain.User{}, errors.New("токен не найден")
	}
	return usr, nil
}

func (s *Secret) Delete(token string) {
	s.mu.Lock()
	delete(s.tokens, token)
	s.mu.Unlock()
}

func Token(r *http.Request) string { //получает токен из HTTP-запроса
	if t := r.Header.Get("Authorization"); t != "" {
		return t
	}
	if c, err := r.Cookie("token"); err == nil {
		return c.Value
	}
	return ""
}
