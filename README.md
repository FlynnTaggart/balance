# Balance Microservice
## Микросервис для работы с балансом пользователей
## Тестовое задание на позицию стажёра-бэкендера AvitoTech

[Задание](https://github.com/avito-tech/internship_backend_2022)

Микросервис предоставляет HTTP API для работы с балансом пользователей. 
Реализованы следующие функции:

* Зачисление средств на баланс пользователя
* Получение баланса пользователя
* Резервирование средств перед совершение покупки
* Совершение покупки (подтверждение списания зарезервированных средств)
* Разрезервирование средств, если покупку совершить не удалось (автоматически и по запросу)
* Формирование CSV отчета о выручке по каждой услуге за расчетный период (месяц)

Дополнительные функции
* Добавление списка услуг
* Получение информации об услуге
* Удаление пользователя 
* Удаление услуги

## Реализация

Микросервис разработан на Go. Для хранения информации используется реляционная 
СУБД - PostgreSQL. Сервис запускается в среде docker-compose один контейнер для сервиса,
и один для базы данных. 

В качестве драйвера для работы с PostgreSQL на Go был выбран [PGX](https://github.com/jackc/pgx), как современное 
и быстрое решение. В нем реализованно множество функций, как и в общем SQL, так и конкретно
PostgreSQL. В нем удобно реализован пул соедениений с базой - [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v4/pgxpool).
Также PostgreSQL позволяет защищать БД от ошибок, возникющих при конкуретном выполенении
транзакций, поэтому, чтобы не терять данные, все транзакции выполнялись в `Serializable` режиме.

В начале разработке использовался [GORM](https://gorm.io/), однако после изучения огромного
количества материалов, посвященных работе с базами данных в Go, было принято решение
отказаться от ORM в пользу скорости и вариативности, но немного в ущерб читаемости 
и краткости кода.

В качестве фреймворка для рутинга был выбран [Fiber](https://github.com/gofiber/fiber).
Он очень похож на другие рутинг фреймворки, но существенно выигрывает по скорости.

В качестве ферймворка для логирования был выбыран [Zap](https://github.com/uber-go/zap).
Он поддерживается PGX (логируются запросы). Логи сохраняются в файл в JSON формате 
(директория с логами смонтирована в docker compose).

Также для генерации swagger файлов был использован [swag](https://github.com/swaggo/swag).

## Запуск

1. Актуализируйте образы контейнеров
```
docker compose pull
```
2. Запустите контейнеры (явно указав, что нужно собрать образы, чтобы точно подтянуть актуальный код)
```
docker compose up -d --build
```

При первом запуске при отправке первого запроса может перестать работать контейнер с сервисом,
так как у него не получается подключиться к БД. Поэтому возможно придется выполнить команду
запуска еще раз (`docker compose up`).

## Что удалось, а что нет

Удалось выполнить основное задание, первое дополнительно задание, удалось реализовать
функцию резрезервироания средств (релизована, с помощью конкурентного запуска таймера), 
удалось сгенерировать `swagger` файл для API.

Не хватило времени, чтобы написать тесты, а также реализовать второе дополнительное задание.
Для выполнения второго дополнительного задания была подготовлена почва в виде таблицы в БД,
хранящей все операции (с ее помощью также реализованы отчеты для бухгалтерии, используются
индексы).

## Запросы

Для того наглядно попробовать все запросы, после запуска сервиса, перейдите в браузере по адресу
`http://localhost:8080/swagger/index.html`. Там будут описаны все запросы, 
также там их можно будет опробовать.
