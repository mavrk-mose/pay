FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

ENV SERVICE=pay_WALLET
ENV ID=pay_WALLLET_1
ENV PORT=8060

RUN go build -o main ./cmd/main.go

EXPOSE 8060

CMD ["./main"]