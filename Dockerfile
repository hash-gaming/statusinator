FROM golang:buster as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test -covermode=atomic -coverpkg=all ./...

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch

COPY --from=builder /app/statusinator /app/
ENTRYPOINT ["/app/statusinator"]
