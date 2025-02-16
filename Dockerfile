FROM golang:1.23.3-alpine3.20

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o ./ cmd/app/main.go

EXPOSE 9090

CMD ["./main"]