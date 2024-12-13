FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

ENV SERVICE=pay_ledger
ENV ID=pay_LEDGER_1
ENV PORT=8080

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]