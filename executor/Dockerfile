FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

ENV SERVICE=pay_executor
ENV ID=pay_EXECUTOR_1
ENV PORT=8090

RUN go build -o main ./cmd/main.go

EXPOSE 8090

CMD ["./main"]