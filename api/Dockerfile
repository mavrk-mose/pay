FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

ENV SERVICE=pay_api
ENV ID=pay_API_1
ENV PORT=8100

RUN go build -o main ./cmd/main.go

EXPOSE 8100

CMD ["./main"]