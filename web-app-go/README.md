```
go mod init project

go get github.com/gin-gonic/gin@latest
go get github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/files@latest
go get github.com/swaggo/gin-swagger@latest

go mod tidy

swag init -g cmd/server/main.go -o docs

go run ./cmd/server

docker-compose up --build

```