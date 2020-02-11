FROM golang:buster as builder

RUN apt-get update && apt-get install ca-certificates

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test -covermode=atomic -coverpkg=all ./...

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/statusinator /app/
ENTRYPOINT ["/app/statusinator"]
