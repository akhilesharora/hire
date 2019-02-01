FROM golang:1.11.1 AS builder

WORKDIR /app
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o /app/messagebird .

FROM alpine:latest

# Support for sqlite
RUN apk --update upgrade
RUN apk add sqlite
RUN apk add curl
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder /app/messagebird ./
COPY --from=builder /app/docker-entrypoint.sh ./
EXPOSE 9365
ENTRYPOINT ["./docker-entrypoint.sh"]
