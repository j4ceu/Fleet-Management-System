FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o fleet ./cmd/main.go
RUN go build -o publisher ./cmd/publisher/main.go
RUN go build -o worker ./cmd/worker/main.go

EXPOSE 8080
CMD ["./fleet"]