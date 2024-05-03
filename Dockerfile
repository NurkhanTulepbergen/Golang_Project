FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o golang-final-project cmd/main.go

CMD ["./golang-final-project"]