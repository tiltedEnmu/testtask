# Тестовое задание

Как запустить:
1. Клонируем себе репозиторий `git clone github/tiltedEnmu/testtask`
2. Поднимаем кафку и постгрес `docker compose up -d` (если образов нет процесс не быстрый)
3. Запускаем продюсер `go run ./producer/cmd/main.go --config=config/config.yaml`
4. И консюмер `go run ./consumer/cmd/main.go --config=config/config.yaml`
5. На `http://localhost:8080/invoice` (или на другой адрес который в конфиге продюсера) отправляем сообщение формата `
   {
   "wallerId": 1,
   "currCode": "RUB",
   "amount": 100.5
   }
`
6. Так же делаем и с withdraw и с balance

P.S. Код очень сырой т.к. было мало времени, но в целом вроде работает.