### Для запуска сервиса необходимо воспользоваться командой:
```bash
make run
```

### Запуск Postgres и nats-streaming:
```bash
docker-compose build
docker-compose up -d
```

### Запуск скрипта для тестовых публикаций:
```bash
make publish
```

### WRK-тест:
```bash
make wrk
```

### Запуск миграций
```bash
make migrate-up
```

### Откат миграций
```bash
make migrate-up
```