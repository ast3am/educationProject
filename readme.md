# educationProject

HTTP-сервис, который принимает входящие соединения с JSON-данными и обрабатывает их следующим образом:

Создание пользователя, пример запроса:
POST /create HTTP/1.1 Content-Type: application/json; charset=utf-8 Host: localhost:8080 {"name":"some name","age":"24","friends":[]}

Данный запрос должен возвращать ID пользователя и статус 201.

Создание друзей, пример запроса:
POST /make_friends HTTP/1.1 Content-Type: application/json; charset=utf-8 Host: localhost:8080 {"source_id":"1","target_id":"2"}

Данный запрос должен возвращать статус 200 и сообщение «username_1 и username_2 теперь друзья».

Удаление пользователя, пример запроса:
DELETE /user HTTP/1.1 Content-Type: application/json; charset=utf-8 Host: localhost:8080 {"target_id":"1"}

Данный запрос должен возвращать 200 и имя удалённого пользователя.

Возвращение всех друзей пользователя:
GET /friends/user_id HTTP/1.1 Host: localhost:8080 Connection: close

Данный запрос должен возвращать 200 и список друзей запрашиваемого пользователя

Обновление возраста пользователя, пример запроса:
PUT /user_id HTTP/1.1 Content-Type: application/json; charset=utf-8 Host: localhost:8080 {"new_age":"28"}

Запрос должен возвращать 200 и сообщение «возраст пользователя успешно обновлён».