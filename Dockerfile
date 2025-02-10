FROM --platform=linux/amd64 golang:1.23-bullseye AS builder

WORKDIR /avito-shop/

COPY . /avito-shop/

RUN go build -o ./api ./cmd/api \
    && go clean -cache -modcache

EXPOSE 8080

ENTRYPOINT ["./api"]
