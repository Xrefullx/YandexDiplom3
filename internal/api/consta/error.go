package consta

import "errors"

const (
	ErrorWorkDataBase    = "ошибка работы с базой данных"
	ErrorUnmarshalBody   = "ошибка Unmarshal тело запроса"
	ErrorReadBody        = "ошибка чтения тело запроса"
	ErrorNumberValidLuhn = "ошибка неверный формат номера заказа"
)

var ErrorNoUNIQUE = errors.New("ошибка значение не уникально")
var ErrorStatusShortfallAccount = errors.New("на счету недостаточно средств")
