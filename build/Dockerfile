FROM golang:1.19.3-alpine3.16 as build

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/app/main.go

FROM alpine:3.6

WORKDIR /app

COPY --from=build /app/main .

ENTRYPOINT ["./main"]
