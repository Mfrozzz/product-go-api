from golang:1.24.3

WORKDIR /go/src/app

COPY . .

EXPOSE 8000

RUN go build -o main cmd/main.go

CMD ["./main"]