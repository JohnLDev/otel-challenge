FROM golang:1.21.0 as builder

WORKDIR /go/app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/app/main -ldflags="-w -s" main.go

FROM scratch as runner

WORKDIR /go/app

COPY --from=builder /go/app/main .
COPY .env .env

ENTRYPOINT [ "./main"]
