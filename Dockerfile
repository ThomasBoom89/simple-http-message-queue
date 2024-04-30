FROM golang:1.22-alpine3.19 AS Builder

WORKDIR /app
COPY . .

RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -o application ./src/main.go

FROM scratch AS Production

USER 1000

WORKDIR /app

COPY --from=Builder /app/application /app

EXPOSE 3000

CMD ["./application"]

FROM golang:1.22-alpine3.19 AS Development

RUN go install github.com/cosmtrek/air@v1.51.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.22.1

WORKDIR /app

CMD ["air"]
