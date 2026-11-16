package middle

import "shop/secret"

type Middle struct {
	secret *secret.Secret //хранилище токенов
}

func New(secret *secret.Secret) *Middle {
	return &Middle{secret: secret} //создаем миддл и запоминаем сикрет
}
