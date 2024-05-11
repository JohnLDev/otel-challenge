FROM golang:1.21.0 as builder

WORKDIR /go/app2

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/app2/main -ldflags="-w -s" main.go

FROM scratch as runner

WORKDIR /go/app2

COPY --from=builder /go/app2/main .
# COPY .env .env

ENTRYPOINT [ "./main"]
