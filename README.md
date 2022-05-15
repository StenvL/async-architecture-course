# async-architecture-course

Zero homework Miro:
https://miro.com/app/board/uXjVO6Eburg=/?moveToWidget=3458764523942276695&cot=10

## Процесс миграции на новую версию событий
**Постановка**:  

1. Добавляем столбец в таблицу с тасками для ключа таска. Назовём его key и будем считать, что префикс всегда будет POPUG. Тогда столбец можно сделать автоинкрементным. 
2. Создаём 2-ую версию события создания таска, в которой добавлено поле key. 
3. Описываем схему 2-ой версии события создания таска в schema registry, там же (в описании схемы) делаем валидацию отсутствия [] в title таска. 
4. Добавляем консьюмер в сервис billing, способный обрабатывать события новой версии, деплоим. 
5. Создаём продьюсер событий новой версии. Валидируем перед отправкой по схеме, описанной в п.3, деплоим.
6. Со временем убеждаемся, что продьюсеров 1-ой версии не осталось, после чего можем выпилить консьюмер 1-ой версии.

## Run
Auth:
* docker-compose build
* docker-compose run oauth rake db:create
* docker-compose run oauth rake db:migrate
* docker-compose up  

Tracker:
* cd tracker
* make migrations.up
* go run main.go
* Создать пользователя по адресу http://localhost:3000, выдать админские права
* Создать application по адресу http://localhost:3000/oauth/applications, указать return URL http://localhost:8080/swagger/oauth2-redirect.html
* Открыть swagger-интерфейс по адресу http://localhost:8080/swagger/index.html, scopes: read write admin
* Авторизоваться и выполнять запросы
  
Billing:
* cd billing
* make migrations.up
* go run main.go
* Создать пользователя по адресу http://localhost:3000, выдать админские права
* Создать application по адресу http://localhost:3000/oauth/applications, указать return URL http://localhost:8081/swagger/oauth2-redirect.html
* Открыть swagger-интерфейс по адресу http://localhost:8081/swagger/index.html, scopes: read write admin
* Авторизоваться и выполнять запросы