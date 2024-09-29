FROM golang:1.22.1-alpine AS base

RUN apk update && apk add --no-cache ca-certificates

FROM base AS builder
WORKDIR /app
COPY . .
RUN go build -o order cmd/order/cmd/main.go
RUN go build -o user cmd/user/cmd/main.go
RUN go build -o product cmd/product/cmd/main.go
RUN go build -o payment cmd/payment/cmd/main.go

FROM alpine:3.19 AS final
WORKDIR /app
COPY --from=builder /app/order .
COPY --from=builder /app/user .
COPY --from=builder /app/product .
COPY --from=builder /app/payment .
COPY .env .

EXPOSE 3330 3332 3333 3331 3430 3431 3432 3433

CMD ["/app/order"]

