FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /main cmd/api/main.go

EXPOSE 9000

ENTRYPOINT ["/main"]
