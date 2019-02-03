FROM golang:1.11.1 AS builder

WORKDIR /app
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o /app/hire .

FROM alpine:latest

RUN apk --update upgrade
RUN apk add curl

COPY --from=builder /app/hire ./
COPY --from=builder /app/docker-entrypoint.sh ./
EXPOSE 8070
ENTRYPOINT ["./docker-entrypoint.sh"]
