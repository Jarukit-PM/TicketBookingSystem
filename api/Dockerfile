FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 go build -o /server ./cmd/server
RUN CGO_ENABLED=0 go build -o /worker ./cmd/worker

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /server /server
COPY --from=builder /worker /worker

EXPOSE 8080

CMD ["/server"]
