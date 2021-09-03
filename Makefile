build:
	go build -o ./bin/qbot cmd/main.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/qbot_linux cmd/main.go

run:
	go run cmd/main.go