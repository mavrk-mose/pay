FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

ENV SERVICE=payment_service
ENV ID=pay_API_1
ENV PORT=8100

RUN go build -o main ./cmd/main.go

EXPOSE 8100

CMD ["./main"]
