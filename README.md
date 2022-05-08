# async-architecture-course

Zero homework Miro:
https://miro.com/app/board/uXjVO6Eburg=/?moveToWidget=3458764523942276695&cot=10

Run:
* docker-compose run oauth rake db:create
* docker-compose run oauth rake db:migrate
* docker-compose up
* cd tracker
* make migrations.up
* go run main.go
* Создать пользователя по адресу http://localhost:3000, выдать админские права
* Создать application по адресу http://localhost:3000/oauth/applications, указать return URL http://localhost:8080/swagger/oauth2-redirect.html
* Открыть swagger-интерфейс по адресу http://localhost:8080/swagger/index.html, scopes: read write admin
* Авторизоваться и выполнять запросы