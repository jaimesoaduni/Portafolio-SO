FROM golang:latest

WORKDIR /app

COPY assets assets
COPY js js
COPY public public
COPY templates templates
COPY assets assets
COPY css css

COPY go.sum .
COPY go.mod .

COPY main.go .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

EXPOSE 3000

CMD ["/app/main"]