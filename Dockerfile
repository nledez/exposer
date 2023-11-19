FROM golang:1.21.4-bookworm as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go .

RUN go build
RUN ls -alrt

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/exposer .

EXPOSE 8080

CMD ["./exposer"]