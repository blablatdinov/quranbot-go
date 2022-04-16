run:
	go run cmd/bot/main.go

test:
	go test -v ./...

lint:
	go fmt ./...