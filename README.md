# Web_calculator

Веб-сервис, вычисляющий простые арифметические выражения.

Обмен происходит в формате JSON. Пользователь отправляет на сервер, на котором запущено приложение (веб-сервис), по адресу /api/v1/calculate POST-запрос с телом вида
`{ "expression": "(1+3) / ((2+0)*2)" }`.
В ответ получает ответ вида 
`{"result":"1.000000"}` 
либо ошибку. Для этого удобно использовать утилиту curl:
```
curl \
-H 'Content-Type: application/json' \
-d '{ "expression": "(1+3) / ((2+0)*2)" }' \
-X POST \
-w "\nHTTP status code %{http_code}" \
'localhost:8877/api/v1/calculate'
```

Примеры запросов успешных и с ошибками:
![см. requests.png](./requests.png)

Имитация внутренней ошибки сервера происходит при отправке любого http-метода, отличного от GET и POST.

Настройки читаются из конфига в формате yaml. При работе приложения в stdout выводится лог, пример:
```
$ ./main.exe
time=2024-12-22T19:49:43.081+03:00 level=INFO msg="starting web_calculator app"
time=2024-12-22T19:49:43.081+03:00 level=INFO msg="server started on 127.0.0.1:8877"
time=2024-12-22T19:50:32.512+03:00 level=INFO msg="connection from 127.0.0.1:58950, method POST"
time=2024-12-22T19:51:54.006+03:00 level=INFO msg="connection from 127.0.0.1:60086, method GET"
time=2024-12-22T19:51:58.281+03:00 level=INFO msg="connection from 127.0.0.1:60356, method DELETE"
time=2024-12-22T19:52:02.577+03:00 level=INFO msg="stopping server"
time=2024-12-22T19:52:02.580+03:00 level=INFO msg="server stopped"
```

## Запуск
1. Установить необходимые пакеты:
`go get github.com/ilyakaznacheev/cleanenv`
2. Собрать проект:
`go build cmd/main.go`
3. При необходимости отредактировать конфиг config.yaml, можно установить его кастомное расположение через переменную среды WCALC_CONFIG_PATH
4. Запустить полученный исполняемый файл.
