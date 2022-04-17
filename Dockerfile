FROM golang:1.18-alpine AS build

WORKDIR /src/
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/qbot cmd/bot/main.go

FROM scratch
COPY --from=build /bin/qbot /bin/qbot
ENTRYPOINT ["/bin/qbot"]