FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . /app

ENV SERVICE=pay_service
ENV ID=pay_SERVICE_1
ENV PORT=8070

RUN go build -o main ./cmd/main.go

EXPOSE 8070

CMD ["./main"]