build:
	go build -o ./bin/qbot cmd/bot/main.go

build_utils:
	go build -o ./bin/qbot_utils cmd/utils/main.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/qbot_linux cmd/bot/main.go

run:
	go run cmd/bot/main.go

lint:
	go fmt .

test:
	go test -v ./...
