package consta

import "errors"

const (
	ErrorDataBase        = "ошибка базы данных"
	ErrorBody            = "ошибка тела запроса"
	ErrorReadBody        = "ошибка чтения запроса"
	ErrorNumberValidLuhn = "неверный формат номера заказа"
)

var ErrorNoUNIQUE = errors.New("значение не уникально")
var ErrorStatusShortfallAccount = errors.New("на счету недостаточно средств")
