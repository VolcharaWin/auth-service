# **AUTHENTICATION API**

## Это _REST API_, которое я сделал для того, чтобы использовать в будущем.

Сейчас оно не полностью закончено, но позволяет запустить на удалённом сервере сервер *gin*, на который, при правильной настройке, можно отправлять запросы и получать ответы от сервера.
На данный момент реализовано только 3 эндпоинта:

1. **/login** - принимает 2 заголовка `login` и `password` и проверяет наличие в базе данных того самого логина и хэша пароля.
2. **/register** - проверяет отсутствие логина в базе данных и записывает новые данные.
3. **/protected/data** - проверка у пользователя куки с аутентификатором, просмотр данных при наличии его.

В будущем планируется улучшить API до современных стандартов.

_Для того, чтобы поставить API на свой сервер, требуется изменить `localhost` в `/server` и в `/authorization`._