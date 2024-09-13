# GophKeeper
Менеджер паролей


### [Техническое задание](./technical%20task.md)

### [Диаграмма взаимодействия](./doc/usecases.plantuml)

### Работа с миграциями БД

Создать новую миграцию
```shell
make migrate-create
```

Применить новые миграции
```shell
make migrate-up
```

Откатить последнюю миграцию
```shell
make migrate-down
```

Откатить все миграции
```shell
make migrate-down-all
```

Внимание! При работе с миграциями сохранность данных не гарантируется и зависит от написанных разработчиком запросов. Прежде чем выполнять то или иное действие - убедись, что ты осознаешь, что ты делаешь.

### Запуск тестов
Перед запуском тестов необходимо создать конфигурации клиента (client.env) на основе файла [client.env.sample](client.env.sample) и сервера (server.env) на основе файла [server.env.sample](server.env.sample).

Требуется установленный docker compose
```shell
make test
```

### Генерация моков
```shell
make build-mocks
```

### Линтер
```shell
make static-check
```

### Форматирование
```shell
make format
```